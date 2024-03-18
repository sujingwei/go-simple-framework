/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 20:01:16
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-15 16:47:57
 * @FilePath: \go-simple-framework\controller\welcomeController.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

type WelcomeController struct {
}

func (controller *WelcomeController) GetHello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "WelcomeController#hello")
}

func (controller *WelcomeController) PostSayHi(ctx *gin.Context) {
	ctx.String(http.StatusOK, "WelcomeController#sayHi")
}

/**
 * @description: 创建cookie操作
 * @param {*gin.Context} ctx
 * @return {*}
 */
func (c *WelcomeController) GetCreateCookie(ctx *gin.Context) {
	webframework.SetCookie(ctx, "userName", "zhangxiaoming", 3600)
	ctx.String(http.StatusOK, "set Cookie Success!")
}

func (c *WelcomeController) GetGetCookie(ctx *gin.Context) {
	if cookie, err := webframework.GetCookie(ctx, "userName"); err == nil {
		ctx.String(http.StatusOK, fmt.Sprintf("get Cookie: %s\n", cookie))
	} else {
		ctx.String(http.StatusOK, fmt.Sprintf("get Cookie err: %+v\n", err))
	}
}

/**
 * @description: 测试xss注入
 * @param {*gin.Context} ctx
 * @return {*}
 */
func (c *WelcomeController) PostXss(ctx *gin.Context) {
	ctx.Request.ParseForm() // 解释表单
	form := make(map[string]any)
	for key, value := range ctx.Request.PostForm {
		form[key] = value[0]
	}
	ctx.JSON(http.StatusOK, gin.H{
		"": form,
	})
}
