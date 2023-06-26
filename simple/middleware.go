package simple

import (
	"net/http"

	pathmatcher "github.com/gohutool/boot4go-pathmatcher"
)

type Middleware struct {
	Rule          string                                                     // 过滤规则，/index/**/*ab
	SkipUrl       []string                                                   // 白名单URL
	middlewareFun func(w http.ResponseWriter, r *http.Request) (bool, error) // 中间件，true，表示继续执行，false表示中止操作
}

// 过滤器集合
var middlewares []*Middleware = make([]*Middleware, 0, 64)

// 添加中间件，包含白名单URL
func UseMiddlewareSkipUrl(rule string, skipUrl []string, fun func(w http.ResponseWriter, r *http.Request) (bool, error)) {
	m := &Middleware{
		Rule:          rule,
		SkipUrl:       skipUrl,
		middlewareFun: fun,
	}
	middlewares = append(middlewares, m)
}

// 添加中间件，不包含白名单URL
func UseMiddleware(rule string, fun func(w http.ResponseWriter, r *http.Request) (bool, error)) {
	UseMiddlewareSkipUrl(rule, nil, fun)
}

// url匹配
// PathMatcher("**/b", "/b")=true
// PathMatcher("**/b", "/b/aaa/a.log")=false
// PathMatcher("**/b/**", "/b/aaa/bbb/a.log")=true
// PathMatcher("**/b/**/", "/b/aaa/bbb/ccc/a.log")=false
// PathMatcher("**/b/**/", "/b/aaa/bbb/ccc/dddd")=false
// PathMatcher("**/b/**", "/b/aaa/bbb/ccc/dddd")=true
// PathMatcher("**/b/**", "/b/aaa/bbb/ccc/dddd")=true
// PathMatcher("**/b/**", "/b/aaa/bbb/ccc/dddd/a.log")=true
// PathMatcher("**/b/**/*.log", "/b/aaa/bbb/a.log")=true`
func PathMatcher(p1, p2 string) (bool, error) {
	return pathmatcher.Match(p1, p2)
}
