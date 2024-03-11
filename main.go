/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 19:48:31
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-11 13:55:15
 * @FilePath: \go-simple-framework\main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sujingwei/go-simple-framework/app/route"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

/**
 * @description:
 * @return {*}
 */
func main() {
	var r *gin.Engine = gin.New()
	route.RetistryRouter(r)
	webframework.WebStart(r)
}
