package article

import (
	logger2 "goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

//Get 通过文章ID 获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}

// 获取全部文章
func GetAll() ([]Article, error)  {
	var articles []Article
	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
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