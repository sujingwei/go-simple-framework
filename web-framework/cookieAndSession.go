/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-13 15:28:23
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-13 15:38:58
 * @FilePath: \go-simple-framework\web-framework\cookieAndSession.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package webframework

import "github.com/gin-gonic/gin"

/**
 * @description: 设置cookie并指定超时时间
 * @param {*gin.Context} ctx
 * @param {*} name
 * @param {string} value
 * @param {int} maxAge
 * @return {*}
 */
func SetCookie(ctx *gin.Context, name, value string, maxAge int) {
	var (
		path     string = "/"
		domain   string = ctx.Request.Host
		secure   bool   = false // 是否智能通过https访问
		httpOnly bool   = true  // 是否允许通过js获取自己的cookie
	)
	ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

/**
 * @description: 获取cookie
 * @param {*gin.Context} ctx
 * @param {string} name
 * @return {*}
 */
func GetCookie(ctx *gin.Context, name string) (val string, err error) {
	val, err = ctx.Cookie(name)
	return
}
