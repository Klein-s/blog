package controller

import (
	"goblog/app/model/article"
	"goblog/app/model/user"
	route "goblog/pkg/routes"
	"goblog/pkg/view"
	"net/http"
)

type UserController struct {
	BaseController
}

//显示用户信息
func (uc UserController) Show(w http.ResponseWriter, r *http.Request)  {

	//获取url id
	id := route.GetRouteVariable("id", r)

	//获取用户信息
	_user, err := user.Get(id)


	if err != nil {
		uc.ResponseForSQLError(w, err)
	} else {

		//读取用户文章
		articles, err := article.GetByUserID(_user.GetStringID())

		//fmt.Fprint(w, articles)
		if err != nil {
			uc.ResponseForSQLError(w, err)
		} else {
			//fmt.Fprint(w, articles)
			view.Render(w, view.D{
				"Article": articles,
			}, "articles.index", "articles._article_meta")
		}
	}


}