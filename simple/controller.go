// 当前包用于配置请求的控制器

package simple

import (
	"errors"
	"strings"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	DELETE string = "DELETE"
	PUT    string = "PUT"
	ANY    string = "ANY"
)

type Controller struct {
	Method        string           // GET、POST、DELETE、PUT、ANY
	RequestUrl    string           // 请求URL，全局唯一，生成ControllerMap的key, 如：/test
	controllerFun func(c *Context) // 控制器方法
}

// 保存控制器的集合
var controllerMap map[string]*Controller = make(map[string]*Controller)

// 添加路由
func __addRoute(method string, url string, fun func(c *Context)) {
	controllerMap[url] = &Controller{
		Method:        method,
		RequestUrl:    url,
		controllerFun: fun,
	}
}
func AddGetRoute(url string, fun func(c *Context)) {
	__addRoute(GET, url, fun)
}
func AddPostRoute(url string, fun func(c *Context)) {
	__addRoute(POST, url, fun)
}
func AddDeleteRoute(url string, fun func(c *Context)) {
	__addRoute(DELETE, url, fun)
}
func AddPutRoute(url string, fun func(c *Context)) {
	__addRoute(PUT, url, fun)
}
func AddAnyRoute(url string, fun func(c *Context)) {
	__addRoute(ANY, url, fun)
}
func AddRoute(method string, url string, fun func(c *Context)) {
	m := strings.ToUpper(method)
	conns := []string{GET, POST, DELETE, PUT, ANY}
	flag := false
	for _, v := range conns {
		if v == m {
			flag = true
			break
		}
	}
	if !flag {
		panic(errors.New("Add Route The param method Fail!"))
	}
	__addRoute(m, url, fun)
}

// 获取路由
func GetRoute(path string) *Controller {
	return controllerMap[path]
}
