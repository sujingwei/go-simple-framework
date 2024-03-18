/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 19:50:35
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-14 16:27:24
 * @FilePath: \go-simple-framework\route\route.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sujingwei/go-simple-framework/app/controller"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

func RetistryRouter(r *gin.Engine) {

	webframework.RetistryController(r, new(controller.WelcomeController))
	webframework.RetistryController(r, new(controller.SessionController))
	// v1 := r.Group("/v1")
	// {
	// 	webframework.RetistryGroupController(v1, new(controller.WelcomeController))
	// }

}
