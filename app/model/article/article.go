package article

import (
	"goblog/app/model"
	route "goblog/pkg/routes"
	"strconv"
)

//Article 文章模型

type Article struct {
	model.BaseModel
	Title string `valid:"title"`
	Body string `valid:"body"`
}

func (a Article) Link() string  {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}