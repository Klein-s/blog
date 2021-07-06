package review

import (
	logger2 "goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	route "goblog/pkg/routes"
	"net/http"
)

//创建评论
func (review *Review) Create() (err error){
	result := model.DB.Create(&review)
	if err := result.Error; err != nil {
		logger2.LogError(err)
	}
	return nil
}

// 获取文章全部评论
func GetArticleReviews(r *http.Request, perPage int, articleId string) ([]Review, pagination.ViewData, error)  {

	//初始化分页实例
	db := model.DB.Where("article_id = ?", articleId).Preload("User").Model(Review{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("articles.show", "id", articleId), perPage)

	//获取视图数据
	viewData := _pager.Paging()

	//获取数据
	var articles []Review
	_pager.Results(&articles)

	return articles, viewData, nil
}
