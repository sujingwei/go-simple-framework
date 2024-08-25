/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-01-26 17:04:32
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2024-02-19 13:23:29
 * @FilePath: \amy-config\amyconfig\amyconfig.go
 * @Description: 环境和配置管理
 * 	1. 读取env, config file, command line 等配置
 *  2. 优化级command line > config file > env
 *  3. 配置文件的优先级，环境配置 > 默认配置
 */
package configuration

import (
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// 添加读写锁
var rwlock sync.RWMutex

/**
 * @description: 读配置, 指定默认目录和文件
 * @return {*}
 */
func LoadFile(c any, path, filename string) {
	if path != "" {
		setDefaultConfigPath(path)
		if filename != "" {
			if index := strings.LastIndex(filename, "."); index >= 1 {
				setDefaultConfigFileName(filename[index:])
				setDefaultConfigFileSuffix([]string{filename[:index+1]})
			} else {
				panic("param filename err!")
			}
		} else {
			panic("param file must!")
		}
	}

	Load(c)
}

/**
 * @description: 读配置
 * @return {*}
 */
func Load(c any) {
	var (
		mStruct    map[string]any = structs.Map(c)
		properties map[string]any = make(map[string]any, 0)
	)
	// struct配置转为properties配置
	mapToProperties(mStruct, properties, "")

	// 从各个环境中读取配置，并保存到properties对象中
	loadEnv(properties)         // 读取 env 配置
	loadFile(properties, c)     // 从配置文件中读取配置
	loadCommandLine(properties) // 从命令行中读取

	// properties配置转为字典配置
	propertiesToMap(mStruct, properties, "")

	// 将map中的数据重新转到c结构体中
	// mapstructure.Decode(mStruct, c)
	mapToStruct(mStruct, c)
}

func mapToStruct(mStruct map[string]any, c any) {
	defer rwlock.Unlock()
	rwlock.Lock()
	mapstructure.Decode(mStruct, c)
}

/**
 * @description: 读取 env 配置
 * @return {*}
 */
func loadEnv(properties map[string]any) {
	for k := range properties {
		if v, b := os.LookupEnv(strings.TrimSpace(strings.ToUpper(strings.ReplaceAll(k, ".", "_")))); b {
			properties[k] = v
		}
	}
}

/**
 * @description: 从配置文件中读取配置
 * @param: properties 转换properties类型map
 */
func loadFile(properties map[string]any, c any) {
	// 获取当前环境配置文件
	var configurationFile *ConfiguratinFileStruct = getConfigurationFile()
	if configurationFile != nil {
		// 读环境配置
		readConfigFile(properties, configurationFile)
	}
	env := getThisEnv(properties, c)
	if env != "" {
		configurationFile.FileName += "-" + env
		configurationFile.FullName = configurationFile.Path + configurationFile.FileName + "." + configurationFile.Suffix
		if fileExists(configurationFile.FullName) {
			readConfigFile(properties, configurationFile)
		}
	}
}

/**
 * @description: 读取 命令行 配置
 * @param: properties
 * @return {*}
 */
func loadCommandLine(properties map[string]any) {
	var (
		args     []string = os.Args   // 命令行入参
		tagStr   string   = "-G"      // 命令行前缀
		splitStr string   = "="       // 分隔符
		argsLen  int      = len(args) // 命令行长度
	)
	for i, v := range args {
		// 获取所有 -G 开始的参数
		if ok := strings.HasPrefix(v, tagStr); ok {
			arg := strings.TrimSpace(strings.ToLower(strings.Replace(v, tagStr, "", 1)))
			if strings.Contains(arg, splitStr) {
				a := strings.Split(arg, splitStr)
				if len(a) == 2 {
					a0, a1 := a[0], strings.TrimSpace(a[1])
					if _, ok := properties[a0]; ok && a1 != "" {
						properties[a0] = a1
					}
				}
			} else {
				if _, ok := properties[arg]; ok && i+1 < argsLen {
					properties[arg] = args[i+1]
				}
			}
		}
	}
}

/**
 * @description: 获取当前环境变量
 */
func getThisEnv(properties map[string]any, c any) string {
	var (
		env string
	)
	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		val := v.Field(i).Interface()
		if _, ok := val.(App); ok {
			if val0, b0 := properties[strings.ToLower(tField.Name)+".env"]; b0 {
				env = val0.(string)
				break
			}
		}
	}
	return strings.TrimSpace(env)
}

/**
 * @description: 从配置文件中读取配置
 * @param: properties 转换properties类型map
 * @param: path
 * @param: filename
 * @param: suffix
 */
func readConfigFile(properties map[string]any, configurationFile *ConfiguratinFileStruct) {
	vip := viper.New()
	vip.AddConfigPath(configurationFile.Path)
	vip.SetConfigName(configurationFile.FileName)
	vip.SetConfigType(configurationFile.Suffix)
	if err := vip.ReadInConfig(); err != nil {
		log.Fatalf("err: %+v\n", err)
		return
	}
	for _, k := range vip.AllKeys() {
		if _, ok := properties[k]; ok {
			v := vip.Get(k)
			if v != nil && v != "" {
				properties[k] = v
			}

		}
	}
}

/**
 * @description: 配置字典转为properties配置
 * @param: mStruct 和配置结构体sttuct对应的map
 * @param: properties 转换properties类型map
 * @param: prefix 配置前缀
 */
func mapToProperties(mStruct map[string]any, properties map[string]any, prefix string) {
	for k, v := range mStruct {
		var kk string
		if prefix == "" {
			kk = k
		} else {
			kk = prefix + "." + k
		}
		if m, ok := v.(map[string]any); ok {
			mapToProperties(m, properties, kk)
		} else {
			// 不区分大小写
			properties[strings.ToLower(kk)] = v
		}
	}
}

/**
 * @description: properties配置转为字典配置
 * @param: mStruct 和配置结构体sttuct对应的map
 * @param: properties 转换properties类型map
 * @param: prefix 配置前缀
 */
func propertiesToMap(mStruct, properties map[string]any, prefix string) {
	for k, v := range mStruct {
		var kk string
		if prefix == "" {
			kk = k
		} else {
			kk = prefix + "." + k
		}
		if m, ok := v.(map[string]any); ok {
			propertiesToMap(m, properties, kk)
		} else {
			// 不区分大小写
			if v2, ok2 := properties[strings.ToLower(kk)]; ok2 && v2 != nil && v2 != "" {
				mStruct[k] = v2
			}
		}
	}
}

/* 获取配置文件 */
func getConfigurationFile() *ConfiguratinFileStruct {
	for _, v := range defaultConfigFileSuffix {
		configFile := defaultConfigPath + defaultConfigFileName + "." + v
		if fileExists(configFile) {
			return &ConfiguratinFileStruct{
				Path:     defaultConfigPath,
				FileName: defaultConfigFileName,
				Suffix:   v,
				FullName: configFile,
			}
		}
	}
	return nil
}

/**
 * @description: 判断文件是否存在
 * @param {string} filename
 * @return {*}
 */
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// log.Fatalf("读配置文件%s失败，异常信息:%s\n", filename, err.Error())
		return false
	}
	return !info.IsDir()
}
