package user

import (
	logger2 "goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/password"
	"goblog/pkg/types"
)

func (user *User) Create() (err error)  {
	if err = model.DB.Create(&user).Error ; err != nil {
		logger2.LogError(err)
		return err
	}
	return nil
}

//根据 id 获取用户信息
func Get(idstr string) (User, error)  {
	var user User
	id := types.StringToInt(idstr)
	if err := model.DB.First(&user, id).Error; err != nil {
		return user, err
	}

	return user, nil

}

//获取三个用户
func Limit3() ([]User, error) {
	var users []User
	if err := model.DB.Limit(3).Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

//根据邮箱获取用户信息
func GetByEmail(email string) (User, error)  {
	var user User
	if err := model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// ComparePassword 对比密码是否匹配
func (u User) ComparePassword(_password string) bool  {
	return password.CheckHas(_password, u.Password)
}
