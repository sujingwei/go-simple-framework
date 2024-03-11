/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 20:01:16
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-11 15:44:29
 * @FilePath: \go-simple-framework\controller\welcomeController.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
