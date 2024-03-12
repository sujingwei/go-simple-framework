/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 12:25:06
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-12 19:40:24
 * @FilePath: \go-simple-framework\web-framework\bootstrap.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package webframework

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	GET     string = "get"
	POST    string = "post"
	DELETE  string = "delete"
	PATCH   string = "patch"
	PUT     string = "put"
	OPTIONS string = "options"
	HEAD    string = "head"
	ANY     string = "any"
)

/**
 * @description: 可以接收的请求类型集合
 * @return {*}
 */
var Methods [8]string = [8]string{GET, POST, DELETE, PATCH, PUT, OPTIONS, HEAD, ANY}

/**
 * @description: 注册对象
 * @return {*}
 */
type retistry struct {
	class   string
	method  string
	handler func(*gin.Context)
}

/**
 * @description: 启动服务
 * @param {*gin.Engine} r
 * @return {*}
 */
func WebStart(r *gin.Engine, webConfig *WebConfig) {
	// var httpServer *http.Server
	httpServer := &http.Server{
		Addr:    ":8001",
		Handler: r,
	}
	if webConfig.Addr != "" {
		httpServer.Addr = webConfig.Addr
	}
	if webConfig.ReadTimeout > 0 {
		httpServer.ReadTimeout = time.Duration(webConfig.ReadTimeout) * time.Second
	}
	if webConfig.ReadHeaderTimeout > 0 {
		httpServer.ReadHeaderTimeout = time.Duration(webConfig.ReadHeaderTimeout) * time.Second
	}
	if webConfig.WriteTimeout > 0 {
		httpServer.WriteTimeout = time.Duration(webConfig.WriteTimeout) * time.Second
	}
	if webConfig.IdleTimeout > 0 {
		httpServer.IdleTimeout = time.Duration(webConfig.IdleTimeout) * time.Second
	}
	if webConfig.MaxHeaderBytes > 0 {
		httpServer.MaxHeaderBytes = webConfig.MaxHeaderBytes
	}
	log.Printf("Start Server: %+v\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("The server[%s] to start failure", httpServer.Addr))
	}
}

/**
 * @description: 将控制器注册到路由组
 * @param {*gin.RouterGroup} g
 * @param {any} controller
 * @return {*}
 */
func RetistryGroupController(g *gin.RouterGroup, controller any) {
	registrys := getRetistryControllerMethod(controller)
	for i := 0; i < len(registrys); i++ {
		re := registrys[i]
		registerGroupControllerRoute(g, re.class, re.method, re.handler)
	}
}

/**
 * @description:  将控制器注册到路由
 * @param {*gin.Engine} r
 * @param {any} controller
 * @return {*}
 */
func RetistryController(r *gin.Engine, controller any) {
	registrys := getRetistryControllerMethod(controller)
	for i := 0; i < len(registrys); i++ {
		re := registrys[i]
		registerControllerRoute(r, re.class, re.method, re.handler)
	}
}

/**
 * @description: 获取需要注册到gin路由的控制器方法
 * @param {any} controller
 * @return {*}
 */
func getRetistryControllerMethod(controller any) []*retistry {
	var rs []*retistry = make([]*retistry, 0)
	if t := reflect.TypeOf(controller); t != nil {
		v := reflect.ValueOf(controller)
		if v.NumMethod() > 0 {
			for i := 0; i < v.NumMethod(); i++ {
				method := v.Method(i)
				// fmt.Printf("method.Name[%s]: %+s\t\t", t.Elem().Name(), t.Method(i).Name)
				// 方法只有1个参数
				if method.Type().NumIn() == 1 {
					// 方法类型为：func(*gin.Context)
					if handler, ok := method.Interface().(func(*gin.Context)); ok {
						// 注册Gin路由
						// _registerControllerRoute(r, t.Elem().Name(), t.Method(i).Name, handler)
						rs = append(rs, &retistry{
							class:   t.Elem().Name(),
							method:  t.Method(i).Name,
							handler: handler,
						})
					}
				}
			}
		}
	}
	return rs
}

/**
 * @description: 注册gin路由
 * @param {*} gin
 * @return {*}
 */
func registerControllerRoute(r *gin.Engine, className, methodName string, handler func(*gin.Context)) {
	for i := 0; i < len(Methods); i++ {
		if strings.HasPrefix(strings.ToLower(methodName), Methods[i]) {
			if Methods[i] == GET {
				r.GET(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), GET, "", 1),
					handler)
			} else if Methods[i] == POST {
				r.POST(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), POST, "", 1),
					handler)
			} else if Methods[i] == DELETE {
				r.DELETE(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), DELETE, "", 1),
					handler)
			} else if Methods[i] == PATCH {
				r.PATCH(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), PATCH, "", 1),
					handler)
			} else if Methods[i] == PUT {
				r.PUT(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), PUT, "", 1),
					handler)
			} else if Methods[i] == OPTIONS {
				r.OPTIONS(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), OPTIONS, "", 1),
					handler)
			} else if Methods[i] == HEAD {
				r.HEAD(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), HEAD, "", 1),
					handler)
			} else {
				r.Any(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), ANY, "", 1),
					handler)
			}
			break
		}
	}
}

/**
 * @description: 注册gin路由
 * @param {*} gin
 * @return {*}
 */
func registerGroupControllerRoute(g *gin.RouterGroup, className, methodName string, handler func(*gin.Context)) {
	for i := 0; i < len(Methods); i++ {
		if strings.HasPrefix(strings.ToLower(methodName), Methods[i]) {
			if Methods[i] == GET {
				g.GET(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), GET, "", 1),
					handler)
			} else if Methods[i] == POST {
				g.POST(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), POST, "", 1),
					handler)
			} else if Methods[i] == DELETE {
				g.DELETE(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), DELETE, "", 1),
					handler)
			} else if Methods[i] == PATCH {
				g.PATCH(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), PATCH, "", 1),
					handler)
			} else if Methods[i] == PUT {
				g.PUT(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), PUT, "", 1),
					handler)
			} else if Methods[i] == OPTIONS {
				g.OPTIONS(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), OPTIONS, "", 1),
					handler)
			} else if Methods[i] == HEAD {
				g.HEAD(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), HEAD, "", 1),
					handler)
			} else {
				g.Any(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), ANY, "", 1),
					handler)
			}
			break
		}
	}
}
