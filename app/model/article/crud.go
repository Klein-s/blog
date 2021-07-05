package article

import (
	logger2 "goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	route "goblog/pkg/routes"
	"goblog/pkg/types"
	"net/http"
)

//Get 通过文章ID 获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)
	if err := model.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}

// 获取全部文章
func GetAll(r *http.Request, perPage int) ([]Article, pagination.ViewData, error)  {

	//初始化分页实例
	db := model.DB.Model(Article{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("articles.index"), perPage)

	//获取视图数据
	viewData := _pager.Paging()

	//获取数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}

//创建文章
func (article *Article) Create() (err error)  {
	result := model.DB.Create(&article)
	if err = result.Error; err != nil {
		logger2.LogError(err)
		return err
	}
	return nil
}

//更新文章
func (article *Article) Update() (rowsAffected int64, err error)  {
	result := model.DB.Save(&article)
	if err = result.Error; err != nil {
		logger2.LogError(err)
		return 0, err
	}
	return result.RowsAffected,nil
}

//删除文章
func (article *Article) Delete() (rowsAffected int64, err error)  {
	result := model.DB.Delete(&article)
	if err = result.Error; err != nil {
		logger2.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

//根据用户id获取全部文章
func GetByUserID(uid string) ([]Article, error)  {
	var article []Article
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Find(&article).Error; err != nil {
		return article, err
	}

	return article,nil
}

//根据分类ID获取文章
func GetByCategoryID(cid string, r *http.Request, perPage int) ([]Article, pagination.ViewData, error)  {

	// 1. 初始化分页实例
	db := model.DB.Model(Article{}).Where("category_id = ?", cid).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("categories.show", "id", cid), perPage)

	// 2. 获取视图数据
	viewData := _pager.Paging()

	// 3. 获取数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}