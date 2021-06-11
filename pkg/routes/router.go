package route

import (
	"github.com/gorilla/mux"
	logger2 "goblog/pkg/logger"
	"net/http"
)

var route *mux.Router

func SetRoute(r *mux.Router)  {
	route = r
}

//Name2URl 通过路由名称获取url
func Name2URL(routeName string, pairs ...string) string  {

	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		logger2.LogError(err)
		return ""
	}

	return url.String()
}

//获取url 参数
func GetRouteVariable(parameterName string, r *http.Request) string  {
	vars := mux.Vars(r)
	return vars[parameterName]
}