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
	"github.com/sujingwei/go-simple-framework/db/nosql"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

/**
 * @description: 默认配置
 * @return {*}
 */
type App struct {
	Name string `json:"name" yaml:"name"`             // 项目名称
	Env  string `default:"dev" json:"env" yaml:"env"` // 当前环境
	// web相关
	Web webframework.WebConfig `json:"web" yaml:"web"` // web相关的配置

	// 数据库存相关
	Mongo nosql.MongoDbConfig `json:"mongo" yaml:"mongo"` // mongoDB连接配置
}

/* 配置文件信息 */
type ConfiguratinFileStruct struct {
	Path     string // 文件路径
	FileName string // 文件名称
	Suffix   string // 文件后缀
	FullName string // 全路径文件名称
}

/**
 * @description: 全局配置文件
 */
var (
	// 默认配置文件路径
	defaultConfigPath string = "./"
	// 默认配置文件名
	defaultConfigFileName string = "bootstrap"
	// 默认配置文件后缀
	defaultConfigFileSuffix []string = []string{"yml", "yaml", "toml", "json"}
)

/* 设置配置文件目录 */
func setDefaultConfigPath(path string) {
	defaultConfigPath = path
}

/* 设置配置文件名 */
func setDefaultConfigFileName(fileName string) {
	defaultConfigFileName = fileName
}

/* 设置配置后缀 */
func setDefaultConfigFileSuffix(suffixs []string) {
	defaultConfigFileSuffix = suffixs
}
