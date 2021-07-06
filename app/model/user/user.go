package user

import (
	"goblog/app/model"
	route "goblog/pkg/routes"
)

//用户模型
type User struct {
	model.BaseModel

	Name string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email string `gorm:"type:varchar(255);default:NULL;unique" valid:"email"`
	Password string `gorm:"type:varchar(255)" valid:"password"`

	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
	Role string `gorm:"type:varchar(255); default:NULL"`
}

func (u User) Link() string  {
	return route.Name2URL("users.show", "id", u.GetStringID())
}

func (u User) HasRole(role string) bool  {
	if u.Role == role{
		return true
	}
	return false
}
