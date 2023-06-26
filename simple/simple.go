package simple

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// web请求全局上下文
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Path    string         //请求路径，如：/index
	Method  string         // 请求类型，GET\POST\PUT\DELETE\ANY
	params  map[string]any // 参数
}

// 创建新的Context
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	c := new(Context)
	c.Writer = w
	c.Request = r
	c.Path = r.URL.Path
	c.Method = strings.ToUpper(r.Method) // 转大写
	return c
}

// Set请求全局参数
func (c Context) Set(key string, value any) {
	c.params[key] = value
}

// 获取请求全局参数
func (c Context) Get(key string) any {
	return c.params[key]
}

// -----------------------------

// url
func (c Context) URL() *url.URL {
	return c.Request.URL
}

// 返回url.Query所有
func (c Context) Query() url.Values {
	return c.URL().Query()
}

// 解释所有url请求参数
func (c Context) QueryMap() map[string]interface{} {
	var result map[string]interface{}
	var values url.Values = c.Query()
	if values != nil {
		result = make(map[string]interface{})
		for key, vals := range values {
			if vals != nil && len(vals) > 0 {
				result[key] = vals[0]
			}
		}
	}
	return result
}

// 解释url请求参数
func (c Context) QueryParam(key string) string {
	return c.Query().Get(key)
}

// 返回post数据
func (c Context) PostFromMap() (result map[string]any) {
	if err := c.Request.ParseForm(); err != nil {
		panic(err)
	}
	for k, v := range c.Request.Form {
		result[k] = v
	}
	return
}

// 获取单个post表单请求参数
func (c Context) PostFormValue(key string) string {
	return c.Request.FormValue(key)
}

// 获取get/post的请求参数
func (c Context) ParamMap() map[string]any {
	if POST == c.Method {
		return c.QueryMap()
	} else if GET == c.Method {
		return c.PostFromMap()
	}
	return nil
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) ValueOf(key string) string {
	if len(c.Request.URL.Query().Get(key)) > 0 {
		return c.Request.URL.Query().Get(key)
	}
	return c.Request.FormValue(key)
}

// 返回字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		log.Fatal("Write error")
	}
}

// 返回json数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		log.Fatal("Write error")
	}
}

func (c *Context) Temp(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		log.Fatal("Write error")
	}
}
