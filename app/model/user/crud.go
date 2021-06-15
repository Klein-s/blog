package user

import (
	logger2 "goblog/pkg/logger"
	"goblog/pkg/model"
)

func (user *User) Create() (err error)  {
	if err = model.DB.Create(&user).Error ; err != nil {
		logger2.LogError(err)
		return err
	}
	return nil
}
