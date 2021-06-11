package article

import (
	route "goblog/pkg/routes"
	"strconv"
)

//Article 文章模型

type Article struct {
	ID int64
	Title, Body string
}

func (a Article) Link() string  {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
}