/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 19:48:31
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-12 19:46:47
 * @FilePath: \go-simple-framework\main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sujingwei/go-simple-framework/app/route"
	"github.com/sujingwei/go-simple-framework/config"
	"github.com/sujingwei/go-simple-framework/configuration"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

var AppConfig config.AppConfigure

func init() {
	// 读配置文件
	configuration.Load(&AppConfig)
}

/**
 * @description:
 * @return {*}
 */
func main() {
	// gin.SetMode("release")
	var r *gin.Engine = gin.Default()
	route.RetistryRouter(r) // 配置 gin 路由
	webframework.WebStart(r, &AppConfig.Web)
}
