/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 12:25:06
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-05-20 16:25:17
 * @FilePath: \go-simple-framework\web-framework\bootstrap.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-10 12:25:06
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-14 18:22:06
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

var (
	// 可以接收的请求类型集合
	Methods [8]string = [8]string{GET, POST, DELETE, PATCH, PUT, OPTIONS, HEAD, ANY}
	// 当前配置的副本
	webConfigCopy WebConfig
)

// 路由注册函数
type RegistryRouteFunc func(*gin.Engine)

/**
 * @description: 注册对象
 * @return {*}
 */
type retistry struct {
	class   string
	method  string
	handler func(*gin.Context)
}

func NewGin(webConfig *WebConfig) *gin.Engine {
	// 将系统配置绑定到当前配置的副本中
	webConfigCopy = *webConfig
	// 通过gin启示web服务
	r := gin.Default()
	// 注册中间件
	registerMiddleware(r)
	return r
}

/**
 * @description: 注册中间件
 * @param {*gin.Engine} r
 * @return {*}
 */
func registerMiddleware(r *gin.Engine) {
	// log.Printf("读取当前配置的副本：%+v\n", webConfigCopy)
	if webConfigCopy.EnableSession {
		// 启用session
		useSessionMiddleware(r)
	}
	if webConfigCopy.Security.Csrf.Enable {
		// 启动csrf
		useCsrfMiddleware(r)
	}
	if webConfigCopy.Security.Xss.Enable {
		//启动 xss
		useXssMiddleware(r)
	}
}

/**
 * @description: 启动服务，当前方法会同步阻塞
 * @param {*gin.Engine} r
 * @return {*}
 */
func WebStart(r *gin.Engine) {
	// var httpServer *http.Server
	httpServer := &http.Server{
		Addr:    ":8001",
		Handler: r,
	}
	// 指定模板
	if webConfigCopy.Template != "" {
		r.LoadHTMLGlob(webConfigCopy.Template)
	}

	if webConfigCopy.Addr != "" {
		httpServer.Addr = webConfigCopy.Addr
	}
	if webConfigCopy.ReadTimeout > 0 {
		httpServer.ReadTimeout = time.Duration(webConfigCopy.ReadTimeout) * time.Second
	}
	if webConfigCopy.ReadHeaderTimeout > 0 {
		httpServer.ReadHeaderTimeout = time.Duration(webConfigCopy.ReadHeaderTimeout) * time.Second
	}
	if webConfigCopy.WriteTimeout > 0 {
		httpServer.WriteTimeout = time.Duration(webConfigCopy.WriteTimeout) * time.Second
	}
	if webConfigCopy.IdleTimeout > 0 {
		httpServer.IdleTimeout = time.Duration(webConfigCopy.IdleTimeout) * time.Second
	}
	if webConfigCopy.MaxHeaderBytes > 0 {
		httpServer.MaxHeaderBytes = webConfigCopy.MaxHeaderBytes
	}
	log.Printf("Start Server: %+v\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("The server[%s] to start failure", httpServer.Addr))
	} else {
		log.Printf("Web Server Start Success, Addr: [%s]\n", httpServer.Addr)
	}
}

// 异步运行web服务
func AsyncWebStart(webConfig *WebConfig, routeFunc RegistryRouteFunc) *gin.Engine {
	var r *gin.Engine = NewGin(webConfig) // 生成Gin
	routeFunc(r)                          // 路由注册
	go WebStart(r)                        // 异步运行gin web服务
	return r
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
				// 方法只有1个参数
				if method.Type().NumIn() == 1 {
					// 方法类型为：func(*gin.Context)
					if handler, ok := method.Interface().(func(*gin.Context)); ok {
						// 当前的方法为gin的路由方法
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
			} else if Methods[i] == ANY {
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
			} else if Methods[i] == ANY {
				g.Any(strings.Replace(strings.ToLower(className), "controller", "", 1)+"/"+strings.Replace(strings.ToLower(methodName), ANY, "", 1),
					handler)
			}
			break
		}
	}
}
