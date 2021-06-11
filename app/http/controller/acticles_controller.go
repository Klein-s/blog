package controller

import (
	"fmt"
	"goblog/app/model/article"
	logger2 "goblog/pkg/logger"
	route "goblog/pkg/routes"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

type ArticlesController struct {

}

/**
文章详情页面
 */
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger2.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示文章
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": types.Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger2.LogError(err)
		tmpl.Execute(w, article)
	}
}

/**
	文章列表
 */
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request)  {
	//执行查询语句，返回结果集
	article, err := article.GetAll()

	if err != nil {
		//数据库错误
		logger2.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		//加载模板
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger2.LogError(err)

		//渲染模板，将所有文章数据传输进去
		tmpl.Execute(w, article)
	}


}