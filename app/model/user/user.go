package user

import "goblog/app/model"

//用户模型
type User struct {
	model.BaseModel

	Name string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email string `gorm:"type:varchar(255);default:NULL;unique" valid:"email"`
	Password string `gorm:"type:varchar(255)" valid:"password"`

	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

func (u User) Link() string  {
	return ""
}
