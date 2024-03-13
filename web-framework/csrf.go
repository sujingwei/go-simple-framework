/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-13 15:48:29
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-13 16:15:00
 * @FilePath: \go-simple-framework\web-framework\csrf.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package webframework

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

/**
 * @description: 使用csrf中间件
 * @param {*gin.Engine} r
 * @return {*}
 */
func useCsrfMiddleware(r *gin.Engine) {
	r.Use(csrfMiddleware())
	r.Use(csrfTokenMiddleware())
}

/**
 * @description: CSRF验证
 * @return {*}
 */
func csrfMiddleware() gin.HandlerFunc {
	csrfMd := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden - CSRF token invalid!"))
		})),
	)
	return adapter.Wrap(csrfMd)
}

/**
 * @description: 在请求头中响应CSRFToken
 * @return {*}
 */
func csrfTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-CSRF-Token", csrf.Token(c.Request))
	}
}
