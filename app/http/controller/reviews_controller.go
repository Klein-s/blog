package controller

import (
	"fmt"
	"goblog/app/model/review"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	route "goblog/pkg/routes"
	"goblog/pkg/types"
	"net/http"
)

type ReviewsController struct {
	BaseController
}

func (rc *ReviewsController) Store (w http.ResponseWriter, r *http.Request)  {


	_review := review.Review{
		UserID: auth.User().ID,
		ArticleID: types.StringToUint64(r.PostFormValue("article_id")),
		Content:  r.PostFormValue("body"),
	}
	errors := requests.ValidateReviewForm(_review)
	if len(errors) == 0 {
		_review.Create()
		if _review.ID > 0 {
			flash.Success("评论成功")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建失败，请联系管理员")
		}

	} else {
		flash.Danger("评论内容不可为空")
	}
	showUrl := route.Name2URL("articles.show", "id", types.Uint64ToString(_review.ArticleID))
	http.Redirect(w, r, showUrl, http.StatusFound)
}
