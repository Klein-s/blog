package category

import (
	"goblog/app/model"
	route "goblog/pkg/routes"
)

type Category struct {
	model.BaseModel
	Name string `gorm:"type:varchar(255); not null;" valid:"name"`
}

func (c Category) Link() string  {
	return route.Name2URL("categories.show", "id", c.GetStringID())
}
