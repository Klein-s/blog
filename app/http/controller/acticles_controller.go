package controller

import (
	"fmt"
	"goblog/app/model/article"
	logger2 "goblog/pkg/logger"
	route "goblog/pkg/routes"
	"goblog/pkg/view"
	"gorm.io/gorm"
	"html/template"
	"net/http"
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


		view.Render(w, "articles.show", article)
	}
}

/**
	文章列表
 */
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request)  {
	//执行查询语句，返回结果集
	articles, err := article.GetAll()

	if err != nil {
		//数据库错误
		logger2.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		//加载模板

		view.Render(w, "articles.index", articles)
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
			fmt.Fprint(w, "插入成功，ID 为"+_article.GetStringID())
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

/**
	文章修改页面
 */
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request)  {
	//获取url参数

	id := route.GetRouteVariable("id", r)

	//读取对应的文章数据
	article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404文章不存在")
		} else {
			//数据库错误
			logger2.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//读取成功 ，渲染表单
		updateUrl := route.Name2URL("articles.update", "id", id)
		data := ArticlesFormData{
			Title: article.Title,
			Body: article.Body,
			URL: updateUrl,
			Errors: nil,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger2.LogError(err)
		tmpl.Execute(w, data)
	}
}

/**
	文章修改
 */
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request)  {
	//获取url参数
	id := route.GetRouteVariable("id", r)

	//获取对应的文章数据
	_article, err := article.Get(id)

	//处理错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章不存在")
		} else {

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//未出现错误

		//验证表单
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)


		if len(errors) == 0 {
			//表单验证通过,  更新数据

			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()

			if err != nil {
				//数据库错误
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
				return
			}

			//更新成功，跳转到文章详情页面

			if rowsAffected > 0 {
				showUrl := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showUrl, http.StatusFound)
			} else {
				fmt.Fprint(w, "未做任何修改")
			}

		} else {
			//表单验证未通过

			updateUrl := route.Name2URL("articles.update", "id", id)
			data := ArticlesFormData{
				Title: title,
				Body: body,
				URL: updateUrl,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			if err != nil {
				logger2.LogError(err)
			}

			tmpl.Execute(w, data)
		}
	}
}

/**
	删除文章
 */
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request)  {
	//获取url参数
	id := route.GetRouteVariable("id", r)

	//获取文章数据
	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			//数据库错误
			logger2.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//未出现错误, 执行删除操作
		rowsAffected, err := _article.Delete()

		//发生错误
		if err != nil {
			//sql 错误
			logger2.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			//发生未知错误
			if rowsAffected > 0 {
				//重定向到文章列表页面
				indexUrl := route.Name2URL("articles.index")
				http.Redirect(w, r, indexUrl, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}

	}
}

