package controller

import (
	"fmt"
	"goblog/app/model/article"
	"goblog/app/model/user"
	logger2 "goblog/pkg/logger"
	route "goblog/pkg/routes"
	"goblog/pkg/view"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {

}

//显示用户信息
func (*UserController) Show(w http.ResponseWriter, r *http.Request)  {

	//获取url id
	id := route.GetRouteVariable("id", r)

	//获取用户信息
	_user, err := user.Get(id)


	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404用户不存在")
		} else {
			logger2.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {

		//读取用户文章
		articles, _ := article.GetByUserID(_user.GetStringID())

		//fmt.Fprint(w, articles)
		if err != nil {
			logger2.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")
		} else {
			//fmt.Fprint(w, articles)
			view.Render(w, view.D{
				"Article": articles,
			}, "articles.index", "articles._article_meta")
		}
	}


}