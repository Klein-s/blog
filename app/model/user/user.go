package user

import "goblog/app/model"

//用户模型
type User struct {
	model.BaseModel

	Name string `gorm:"column:name;type:varchar(255);not null;unique"`
	Email string `gorm:"column:email;type:varchar(255);default:NULL;unique"`
	Password string `gorm:"column:password;type:varchar(255)"`
}
