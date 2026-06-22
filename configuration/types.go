/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-01-26 17:04:32
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-05-20 14:35:58
 * @FilePath: \amy-config\configuration\amyconfig.go
 * @Description: 默认封装配置类
 */
package configuration

import (
	"sync"

	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

// App 默认配置结构体
type App struct {
	Name string `default:"Gin Application" json:"name" yaml:"name"` // 项目名称
	Env  string `default:"dev" json:"env" yaml:"env"`               // 当前环境
	// web相关
	Web webframework.WebConfig `json:"web" yaml:"web"` // web相关的配置

	// 数据库存相关
	// Mongo nosql.MongoDbConfig `json:"mongo" yaml:"mongo"` // mongoDB连接配置
}

// ConfigurationFileStruct 配置文件信息
type ConfigurationFileStruct struct {
	Path     string // 文件路径
	FileName string // 文件名称
	Suffix   string // 文件后缀
	FullName string // 全路径文件名称
}

// 全局配置变量（读写保护）
var (
	configMu sync.RWMutex

	// 默认配置文件路径
	defaultConfigPath string = "./"
	// 默认配置文件名
	defaultConfigFileName string = "bootstrap"
	// 默认配置文件后缀（优先级顺序）
	defaultConfigFileSuffix []string = []string{"yml", "yaml", "toml", "json"}
)

// 命令行参数前缀和分隔符（可通过 SetCommandLinePrefix/SetCommandLineSplit 自定义）
var (
	cmdLinePrefix string = "-G"
	cmdLineSplit  string = "="
)

// #################### Setter 函数（导出，并发安全） ####################

// SetDefaultConfigPath 设置配置文件目录
func SetDefaultConfigPath(path string) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultConfigPath = path
}

// SetDefaultConfigFileName 设置配置文件名
func SetDefaultConfigFileName(fileName string) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultConfigFileName = fileName
}

// SetDefaultConfigFileSuffix 设置配置后缀（优先级顺序）
func SetDefaultConfigFileSuffix(suffixes []string) {
	configMu.Lock()
	defer configMu.Unlock()
	defaultConfigFileSuffix = suffixes
}

// SetCommandLinePrefix 设置命令行参数前缀，默认 "-G"
func SetCommandLinePrefix(prefix string) {
	configMu.Lock()
	defer configMu.Unlock()
	cmdLinePrefix = prefix
}

// SetCommandLineSplit 设置命令行参数分隔符，默认 "="
func SetCommandLineSplit(split string) {
	configMu.Lock()
	defer configMu.Unlock()
	cmdLineSplit = split
}

// #################### 内部 getter（需在持有锁时调用） ####################

func getDefaultConfigPath() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultConfigPath
}

func getDefaultConfigFileName() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultConfigFileName
}

func getDefaultConfigFileSuffix() []string {
	configMu.RLock()
	defer configMu.RUnlock()
	return defaultConfigFileSuffix
}

func getCmdLinePrefix() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return cmdLinePrefix
}

func getCmdLineSplit() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return cmdLineSplit
}

// 保留旧的未导出 setter 以保持向后兼容（内部使用）
func setDefaultConfigPath(path string) {
	SetDefaultConfigPath(path)
}

func setDefaultConfigFileName(fileName string) {
	SetDefaultConfigFileName(fileName)
}

func setDefaultConfigFileSuffix(suffixes []string) {
	SetDefaultConfigFileSuffix(suffixes)
}
