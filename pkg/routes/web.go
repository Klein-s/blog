package routes

import (
	"github.com/gorilla/mux"
	"goblog/app/http/controller"
	"net/http"
)

// 注册网页相关路由
func RegisterWebRoutes(r *mux.Router)  {

	//静态页面
	pc := new(controller.PageController)
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
}