package route

import "github.com/gorilla/mux"

//router 路由对象
var Router *mux.Router

//initialize 初始化路由
func Initialize()  {
	Router = mux.NewRouter()
}

//Name2URl 通过路由名称获取url
func Name2URL(routeName string, pairs ...string) string  {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		//checkError(err)
		return ""
	}

	return url.String()
}