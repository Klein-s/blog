package category

import "goblog/app/model"

type Category struct {
	model.BaseModel
	Name string `gorm:"type:varchar(255); not null;" valid:"name"`
}
