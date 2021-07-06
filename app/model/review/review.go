package review

import (
	"goblog/app/model"
	"goblog/app/model/user"
)

type Review struct {
	model.BaseModel
	UserID uint64 `gorm:"int(11)"`
	User user.User
	ArticleID uint64 `gorm:"int(11)"`
	Content string `gorm:"varchar(255); default:NULL" valid:"content"`
}