package simple

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	store            = sessions.NewFilesystemStore("./", securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
	sessionId string = "go-simple"
)

// 设置session

func (c Context) SetSession(key string, value any) {
	session, _ := store.Get(c.Request, sessionId)
	session.Values[key] = value
	if err := session.Save(c.Request, c.Writer); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// 获取session
func (c Context) GetSession(key string) any {
	session, _ := store.Get(c.Request, sessionId)
	return session.Values[key]
}

// 删除
func (c Context) DelSession(key string) {
	session, _ := store.Get(c.Request, sessionId)
	delete(session.Values, key)
}

// 设置cookie
func (c Context) SetCookie(key, value string, second int) {
	if second <= 0 {
		second = 3600 * 8
	}
	co := http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   second, // 单位s
		HttpOnly: true,
	}
	c.Writer.Header().Add("Set-Cookie", co.String())
}

// 获取cookie
func (c Context) GetCookie(key string) string {
	co, err := c.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return co.Value
}

// 删除cookie

func (c Context) DelCookie(key string) {
	co := http.Cookie{
		Name:     key,
		Value:    "",
		MaxAge:   -1, // 删除cookie
		HttpOnly: true,
	}
	c.Writer.Header().Add("Set-Cookie", co.String())
}
