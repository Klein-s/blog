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
	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}

// 获取全部文章
func GetAll() ([]Article, error)  {
	var articles []Article
	if err := model.DB.Debug().Preload("User").Find(&articles).Error; err != nil {
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