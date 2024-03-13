/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 20:01:16
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-13 15:45:19
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

func (controller *WelcomeController) POSTAbc() {

}

func (controller *WelcomeController) POSTEfg(a, b int) {

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
