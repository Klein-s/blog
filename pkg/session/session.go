package session

import (
	"github.com/gorilla/sessions"
	logger2 "goblog/pkg/logger"
	"net/http"
)

// Store gorilla session 的存储库
var Store = sessions.NewCookieStore([]byte("33446a9dcf9ea060a0a6532b166da32f304af0de"))

// session 当前会话

var Session *sessions.Session

// request 获取会话

var Request *http.Request

// response 写入会话

var Response http.ResponseWriter

// startSession 初始化会话。 中间件当中调用
func StartSession(w http.ResponseWriter, r *http.Request)  {
	var err error

	//store.get 第二个参数是cookie的名称
	//gorilla/session 支持多会话
	Session, err = Store.Get(r, "blog-session")
	logger2.LogError(err)

	Request = r
	Response = w
}

//put 写入键值对应的会话数据
func Put(key string, value interface{}) {
	Session.Values[key] = value
	Save()
}

//get 获取会话数据， 获取数据时请做类型检测

func Get(key string) interface{}  {
	return Session.Values[key]
}

// forget 删除某个会话项
func Forget(key string)  {
	delete(Session.Values, key)
	Save()
}

// flush 删除当前会话
func Flush()  {
	Session.Options.MaxAge = -1
	Save()
}

// save 保持当前会话

func Save()  {

	//非Http 的链接无法使用 secure 和 httponly, 浏览器会报错
	// session.options.secure = true
	// session.options.httponly = true
	err := Session.Save(Request, Response)
	logger2.LogError(err)
}