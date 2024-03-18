/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-14 16:26:12
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-14 17:29:32
 * @FilePath: \go-simple-framework\app\controller\SessionController.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controller

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

type SessionController struct {
}

func (controller *SessionController) GetHello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "WelcomeController#hello")
}

func (controller *SessionController) GetCreateSession(c *gin.Context) {

	if err := webframework.SetSession(c, "sessionAge", rand.Uint32()%20); err == nil {
		c.String(http.StatusOK, "创建session成功")
	} else {
		c.String(http.StatusOK, "创建session失败！")
	}
}

func (controller *SessionController) GetGetSession(c *gin.Context) {
	sessionAge := webframework.GetSession(c, "sessionAge")
	c.String(http.StatusOK, fmt.Sprintf("获取Session, Type:%T, Val:%v", sessionAge, sessionAge))
}
