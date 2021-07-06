package routes

import (
	"github.com/gorilla/mux"
	"goblog/app/http/controller"
	"goblog/app/http/middlewares"
	"net/http"
)

// 注册网页相关路由
func RegisterWebRoutes(r *mux.Router)  {

	//静态页面
	pc := new(controller.PageController)
	//r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	//文章相关页面
	ac := new(controller.ArticlesController)

	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/", ac.Index).
		Methods("GET").Name("articles.index")
	r.HandleFunc("/articles", ac.Store).
		Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create", middlewares.Auth(ac.Create)).
		Methods("GET").Name("articles.create")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", middlewares.Auth(ac.Edit)).
		Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(ac.Update)).
		Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", middlewares.Auth(ac.Delete)).
		Methods("POST").Name("articles.delete")

	//文章分类

	cc := new(controller.CategoriesController)
	r.HandleFunc("/categories/create",  middlewares.Auth(cc.Create)).
		Methods("GET").Name("categories.create")
	r.HandleFunc("/categories/store",  middlewares.Auth(cc.Store)).
		Methods("POST").Name("categories.store")
	r.HandleFunc("/categories/{id:[0-9]+}", cc.Show).Methods("GET").Name("categories.show")

	//文章评论

	rc := new(controller.ReviewsController)

	r.HandleFunc("/reviews/store", middlewares.Auth(rc.Store)).Methods("POST").Name("reviews.store")

	//用户认证

	auc := new(controller.AuthController)

	r.HandleFunc("/auth/register", middlewares.Guest(auc.Register)).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/register", middlewares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")

	r.HandleFunc("/auth/login", middlewares.Guest(auc.Login)).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/dologin", middlewares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")
	r.HandleFunc("/auth/logout", middlewares.Auth(auc.Logout)).Methods("POST").Name("auth.logout")

	//用户相关页面
	uc := new(controller.UserController)

	r.HandleFunc("/users/{id:[0-9]+}", uc.Show).Methods("GET").Name("users.show")

	//静态资源
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	//注册中间件

	//开启会话
	r.Use(middlewares.StartSession)


}