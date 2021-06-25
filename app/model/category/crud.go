package category

import (
	logger2 "goblog/pkg/logger"
	"goblog/pkg/model"
)

func (category *Category) Create() (err error)  {
	if err = model.DB.Create(&category).Error; err != nil {
		logger2.LogError(err)
		return err
	}
	return nil
}

func All() ([]Category, error) {
	var categories []Category
	if err := model.DB.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}