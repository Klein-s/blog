package article

import (
	"goblog/app/model"
	"goblog/app/model/category"
	"goblog/app/model/user"
	route "goblog/pkg/routes"
	"strconv"
)

//Article 文章模型

type Article struct {
	model.BaseModel
	Title string `valid:"title"`
	Body string `valid:"body"`
	UserID uint64 `gorm:"not null;index"`
	User user.User
	CategoryID uint64 `gorm:"not null; default:2; index;" valid:"category_id"`
	Category category.Category
}

func (a Article) Link() string  {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}