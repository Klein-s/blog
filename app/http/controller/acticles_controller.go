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
	"strconv"
	"unicode/utf8"
)

type ArticlesController struct {

}

//创建文章表单数据
type ArticlesFormData struct {
	Title, Body string
	URL string
	Errors map[string]string
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

/**
	文章创建页面
 */
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request)  {

	storeUrl := route.Name2URL("articles.store")
	data := ArticlesFormData{
		Title : "",
		Body : "",
		URL : storeUrl,
		Errors : nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
}

/**
	创建文章
 */
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request)  {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)


	//检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body: body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建失败，请联系管理员")
		}

	} else {

		storeURl := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title: title,
			Body: body,
			URL: storeURl,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

/**
表单验证
*/
func validateArticleFormData(title string, body string) map[string]string {

	errors := make(map[string]string)

	//验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度为3-40字符"
	}

	//验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容必须大于或等于10个字符"
	}
	return errors
}

