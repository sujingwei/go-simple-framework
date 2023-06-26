package simple

import (
	"fmt"
	"log"
	"net/http"
)

// 启动类
func Bootstrap() {
	autoConfig()          // 加载配置
	bootstrapHandleFunc() // 全局handler
	http.ListenAndServe(":8081", nil)
}

// 加载配置
func autoConfig() {

}

// 全局handler
func bootstrapHandleFunc() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle全局异常处理
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%+v\n", err) // 是使用全局异常处理，还是打开异常，可以通过配置文件配置
				w.Write([]byte(fmt.Sprintf("%s", err)))
			}
		}()
		// 创建&Context实例
		c := NewContext(w, r)
		// 如果有中间件，执行中间件
		b, err := execMiddleware(c)
		if err != nil {
			panic(err)
		}
		if b {
			// 中间件执行成功，可以运行控制器
			controller := GetRoute(c.Path)
			if controller == nil {
				panic("Can't Find Route Func! Page 404")
			}
			controller.controllerFun(c) // 执行路由方法
		}
		c.String(200, "success")
	})
}

// 执行中间件，true:所在中间件均执行成功，否则，执行失败
func execMiddleware(c *Context) (rs bool, err error) {
	rs = true                                 // 中间件执行，默认全部通过
	middlewares := useRouteMiddleware(c.Path) // 获取要执行的中间件列表
	if middlewares != nil && len(middlewares) > 0 {
		// 执行中间件
		for _, m := range middlewares {
			rs, err := m.middlewareFun(c.Writer, c.Request)
			if err != nil {
				rs = false // 有异常，设置为false
				panic(err)
			}
			if !rs {
				break
			}
		}
	}
	return rs, err
}

// 获取当前路由要执行的中间件
func useRouteMiddleware(path string) []*Middleware {
	var useMiddlewares []*Middleware
	log.Printf("%+v\n", "aaa")
	if middlewares != nil && len(middlewares) > 0 {
		for _, m := range middlewares {
			if m.SkipUrl != nil && len(m.SkipUrl) > 0 {
				inSkipUrl := false
				for _, skipUrl := range m.SkipUrl {
					if path == skipUrl {
						inSkipUrl = true
						break
					}
				}
				if inSkipUrl {
					// 是白名称URL
					continue
				}
			}
			if m.Rule == "" {
				// 空表示直接添加
				useMiddlewares = append(useMiddlewares, m)
			} else {
				b, err := PathMatcher(m.Rule, path)
				if err != nil {
					panic(err)
				}
				if b {
					if useMiddlewares == nil {
						useMiddlewares = make([]*Middleware, 10, 10)
					}
					useMiddlewares = append(useMiddlewares, m)
				}
			}
		}
	}
	log.Printf("%+v\n", "bbb")
	return useMiddlewares
}
