/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-13 15:28:23
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-14 18:29:14
 * @FilePath: \go-simple-framework\web-framework\cookieAndSession.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package webframework

import (
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func useSessionMiddleware(r *gin.Engine) {
	log.Println("Use Session Middleware!")
	var sessionStore = cookie.NewStore([]byte(strings.ReplaceAll(uuid.New().String(), "-", "")))
	r.Use(sessions.Sessions("my-session", sessionStore))
}

/**
 * @description: 创建session
 * @param {*gin.Context} ctx
 * @param {string} name
 * @param {any} value
 * @return {*}
 */
func SetSession(ctx *gin.Context, name string, value any) error {
	session := sessions.Default(ctx)
	if v2 := session.Get(name); v2 != nil {
		session.Delete(name)
	}
	session.Set(name, value)
	return session.Save()
}

/**
 * @description: 删除session
 * @param {*gin.Context} ctx
 * @param {*} name
 * @return {*}
 */
func DelSession(ctx *gin.Context, name string) {
	session := sessions.Default(ctx)
	if v2 := session.Get(name); v2 != nil {
		session.Delete(name)
	}
}

/**
 * @description: 获取session
 * @param {*gin.Context} ctx
 * @param {string} name
 * @return {*}
 */
func GetSession(ctx *gin.Context, name string) any {
	session := sessions.Default(ctx)
	return session.Get(name)
}

/**
 * @description: 创建cookie并指定超时时间
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
func GetCookie(ctx *gin.Context, name string) (string, error) {
	return ctx.Cookie(name)
}

/**
 * @description: 删除cookie
 * @param {*gin.Context} ctx
 * @param {string} name
 * @return {*}
 */
func DelCookie(ctx *gin.Context, name string) {
	SetCookie(ctx, name, "", 0)
}
