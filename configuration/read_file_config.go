/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-01-26 17:04:32
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2024-02-19 13:23:29
 * @FilePath: \amy-config\amyconfig\amyconfig.go
 * @Description: 环境和配置管理
 *  1. 读取 env, config file, command line 等配置
 *  2. 优先级 command line > config file(环境配置 > 默认配置) > env
 *  3. 配置文件的优先级，环境配置 > 默认配置
 */
package configuration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// LoadFile 读配置，可指定配置文件目录和文件名
// 参数 path 和 filename 均可单独为空，为空时使用默认值
func LoadFile(c any, path, filename string) error {
	if path != "" {
		SetDefaultConfigPath(path)
	}
	if filename != "" {
		// 解析文件名和后缀： "bootstrap.yml" → 文件名="bootstrap", 后缀="yml"
		index := strings.LastIndex(filename, ".")
		if index < 1 || index == len(filename)-1 {
			return fmt.Errorf("param filename[%s] format error, expected like 'bootstrap.yml'", filename)
		}
		SetDefaultConfigFileName(filename[:index])        // "bootstrap"
		SetDefaultConfigFileSuffix([]string{filename[index+1:]}) // "yml"
	}

	return Load(c)
}

// Load 读取配置：struct 默认值 → env → file → command line
// 优先级：command line > file(环境配置 > 默认配置) > env > struct default
func Load(c any) error {
	// 1. 应用 struct 中的 default tag 作为初始值
	if err := applyDefaultTagValues(c); err != nil {
		return fmt.Errorf("apply default tag values: %w", err)
	}

	// 2. struct 展开为 properties map
	mStruct := structs.Map(c)
	properties := make(map[string]any, len(mStruct)*2)
	mapToProperties(mStruct, properties, "")

	// 3. 按优先级从低到高加载配置
	if err := loadEnv(properties); err != nil {
		return fmt.Errorf("load env: %w", err)
	}
	if err := loadFile(properties, c); err != nil {
		return fmt.Errorf("load file: %w", err)
	}
	loadCommandLine(properties)

	// 4. properties 回写 mStruct
	propertiesToMap(mStruct, properties, "")

	// 5. mStruct → struct
	mapToStruct(mStruct, c)

	return nil
}

// applyDefaultTagValues 通过反射读取 struct 的 default tag，设置零值字段的默认值
// 递归处理嵌套 struct 字段
func applyDefaultTagValues(c any) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("applyDefaultTagValues: c must be a non-nil pointer to struct")
	}
	return applyDefaultsRecursive(v.Elem())
}

func applyDefaultsRecursive(v reflect.Value) error {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := v.Field(i)

		// 递归处理嵌套 struct / struct 指针
		if fv.Kind() == reflect.Struct {
			_ = applyDefaultsRecursive(fv)
		} else if fv.Kind() == reflect.Ptr && !fv.IsNil() && fv.Elem().Kind() == reflect.Struct {
			_ = applyDefaultsRecursive(fv.Elem())
		}

		defaultVal := field.Tag.Get("default")
		if defaultVal == "" {
			continue
		}
		if !fv.CanSet() {
			continue
		}
		// 仅当字段为零值时设置默认值
		if !fv.IsZero() {
			continue
		}
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(defaultVal)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if n, err := parseInt(defaultVal); err == nil {
				fv.SetInt(n)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if n, err := parseUint(defaultVal); err == nil {
				fv.SetUint(n)
			}
		case reflect.Bool:
			if b, err := parseBool(defaultVal); err == nil {
				fv.SetBool(b)
			}
		case reflect.Float32, reflect.Float64:
			if f, err := parseFloat(defaultVal); err == nil {
				fv.SetFloat(f)
			}
		}
	}
	return nil
}

// mapToStruct 将 map 解码到 struct（并发安全）
func mapToStruct(mStruct map[string]any, c any) {
	// mapstructure.Decode 是纯函数，不需要锁保护
	_ = mapstructure.Decode(mStruct, c)
}

// loadEnv 读取环境变量配置（优先级最低）
func loadEnv(properties map[string]any) error {
	for k := range properties {
		envKey := strings.ToUpper(strings.ReplaceAll(k, ".", "_"))
		if v, found := os.LookupEnv(envKey); found {
			properties[k] = strings.TrimSpace(v)
		}
	}
	return nil
}

// loadFile 从配置文件中读取配置
func loadFile(properties map[string]any, c any) error {
	// 读取默认配置文件
	configFile := getConfigurationFile()
	if configFile != nil {
		if err := readConfigFile(properties, configFile); err != nil {
			return err
		}
	}

	// 读取环境特定配置（如 bootstrap-dev.yml）
	env := getThisEnv(properties, c)
	if env != "" {
		configMu.RLock()
		suffix := getEnvConfigFileSuffix(configFile)
		fileName := defaultConfigFileName + "-" + env
		envFile := &ConfigurationFileStruct{
			Path:     defaultConfigPath,
			FileName: fileName,
			Suffix:   suffix,
			FullName: filepath.Join(defaultConfigPath, fileName+"."+suffix),
		}
		configMu.RUnlock()

		if fileExists(envFile.FullName) {
			if err := readConfigFile(properties, envFile); err != nil {
				return err
			}
		}
	}
	return nil
}

// loadCommandLine 读取命令行配置（优先级最高）
func loadCommandLine(properties map[string]any) {
	args := os.Args
	prefix := getCmdLinePrefix()
	split := getCmdLineSplit()
	argsLen := len(args)

	for i, v := range args {
		if !strings.HasPrefix(v, prefix) {
			continue
		}
		// 去掉前缀
		arg := strings.TrimSpace(strings.ToLower(strings.Replace(v, prefix, "", 1)))

		if strings.Contains(arg, split) {
			// 格式: -Gkey=value
			a := strings.Split(arg, split)
			if len(a) == 2 {
				key, val := a[0], strings.TrimSpace(a[1])
				if _, ok := properties[key]; ok && val != "" {
					properties[key] = val
				} else if !ok {
					log.Printf("[configuration] unknown key from command line: %s\n", key)
				}
			}
		} else {
			// 格式: -Gkey value
			if _, ok := properties[arg]; ok && i+1 < argsLen {
				properties[arg] = args[i+1]
			} else if !ok {
				log.Printf("[configuration] unknown key from command line: %s\n", arg)
			}
		}
	}
}

// getThisEnv 从 properties 或 struct 反射中获取当前环境名
func getThisEnv(properties map[string]any, c any) string {
	// 快速路径：直接查找 "app.env"（最常见情况）
	if val, ok := properties["app.env"]; ok {
		if s, ok := val.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}

	// 通用路径：反射遍历 struct，查找类型为 App 的字段
	t := reflect.TypeOf(c)
	if t != nil {
		t = t.Elem()
	}
	v := reflect.ValueOf(c)
	if v.IsValid() {
		v = v.Elem()
	}
	if t == nil || !v.IsValid() {
		return ""
	}

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		fieldVal := v.Field(i)
		if !fieldVal.IsValid() {
			continue
		}
		if _, ok := fieldVal.Interface().(App); ok {
			key := strings.ToLower(tField.Name) + ".env"
			if val, exists := properties[key]; exists {
				if s, ok := val.(string); ok {
					return strings.TrimSpace(s)
				}
			}
		}
	}
	return ""
}

// readConfigFile 使用 viper 读取配置文件并将值填入 properties
func readConfigFile(properties map[string]any, configFile *ConfigurationFileStruct) error {
	vip := viper.New()
	vip.AddConfigPath(configFile.Path)
	vip.SetConfigName(configFile.FileName)
	vip.SetConfigType(configFile.Suffix)
	if err := vip.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file %s: %w", configFile.FullName, err)
	}
	for _, k := range vip.AllKeys() {
		if _, ok := properties[k]; !ok {
			continue
		}
		v := vip.Get(k)
		if v == nil {
			continue
		}
		// 只覆盖非零值，保留已有配置
		if isNonZeroValue(v) {
			properties[k] = v
		}
	}
	return nil
}

// isNonZeroValue 判断值是否为非零值（对字符串也排除空串）
func isNonZeroValue(v any) bool {
	if v == nil {
		return false
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.String() != ""
	case reflect.Bool:
		return true // bool 的 false 也是有效值
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true // 0 也是有效的整数值
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() > 0
	default:
		return !rv.IsZero()
	}
}

// mapToProperties 将嵌套 map 展平为 dot 分隔的 properties
func mapToProperties(mStruct map[string]any, properties map[string]any, prefix string) {
	for k, v := range mStruct {
		kk := k
		if prefix != "" {
			kk = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			mapToProperties(nested, properties, kk)
		} else {
			properties[strings.ToLower(kk)] = v
		}
	}
}

// propertiesToMap 将 properties 中的值回填到嵌套 map 中
func propertiesToMap(mStruct, properties map[string]any, prefix string) {
	for k, v := range mStruct {
		kk := k
		if prefix != "" {
			kk = prefix + "." + k
		}
		if nested, ok := v.(map[string]any); ok {
			propertiesToMap(nested, properties, kk)
		} else {
			if v2, ok := properties[strings.ToLower(kk)]; ok && isNonZeroValue(v2) {
				mStruct[k] = v2
			}
		}
	}
}

// getEnvConfigFileSuffix 获取环境特定配置文件的后缀
// 如果找到了默认配置文件则复用其后缀，否则回退到第一个默认后缀
func getEnvConfigFileSuffix(defaultFile *ConfigurationFileStruct) string {
	if defaultFile != nil {
		return defaultFile.Suffix
	}
	if len(defaultConfigFileSuffix) > 0 {
		return defaultConfigFileSuffix[0]
	}
	return "yml"
}

// getConfigurationFile 获取配置文件信息（按后缀优先级查找）
func getConfigurationFile() *ConfigurationFileStruct {
	configMu.RLock()
	defer configMu.RUnlock()

	for _, suffix := range defaultConfigFileSuffix {
		configFile := filepath.Join(defaultConfigPath, defaultConfigFileName+"."+suffix)
		if fileExists(configFile) {
			return &ConfigurationFileStruct{
				Path:     defaultConfigPath,
				FileName: defaultConfigFileName,
				Suffix:   suffix,
				FullName: configFile,
			}
		}
	}
	return nil
}

// fileExists 判断文件是否存在（非目录）
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// #################### 辅助解析函数（用于 default tag） ####################

func parseInt(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n, err
}

func parseUint(s string) (uint64, error) {
	var n uint64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n, err
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	}
	return false, fmt.Errorf("invalid bool value: %s", s)
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%f", &f)
	return f, err
}
