package controller

import (
	"fmt"
	"goblog/app/model/article"
	"goblog/app/requests"
	logger2 "goblog/pkg/logger"
	route "goblog/pkg/routes"
	"goblog/pkg/view"
	"gorm.io/gorm"
	"net/http"
)

type ArticlesController struct {

}

//创建文章表单数据
type ArticlesFormData struct {
	Title, Body string
	Article article.Article
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


		view.Render(w, view.D{
			"Article": article,
		}, "articles.show")
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

		view.Render(w,  view.D{
			"Article": articles,
		}, "articles.index")
	}


}

/**
	文章创建页面
 */
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request)  {

	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

/**
	创建文章
 */
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request)  {


	_article := article.Article{
		Title: r.PostFormValue("title"),
		Body: r.PostFormValue("body"),
	}
	errors := requests.ValidateArticleForm(_article)


	//检查是否有错误
	if len(errors) == 0 {

		_article.Create()
		if _article.ID > 0 {
			//fmt.Fprint(w, "插入成功，ID 为"+_article.GetStringID())
			showUrl := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, showUrl, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建失败，请联系管理员")
		}

	} else {

		view.Render(w, view.D{
			"Article": _article,
			"Errors": errors,
		}, "articles.create", "articles._form_field")
	}
}

/**
	文章修改页面
 */
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request)  {
	//获取url参数

	id := route.GetRouteVariable("id", r)

	//读取对应的文章数据
	_article, err := article.Get(id)

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

		view.Render(w, view.D{
			"Article": _article,
			"Errors": map[string]string{"title":"", "body":""},
		}, "articles.edit", "articles._form_field")
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
		_article.Title = r.PostFormValue("title")
		_article.Body = r.PostFormValue("body")

		errors := requests.ValidateArticleForm(_article)


		if len(errors) == 0 {
			//表单验证通过,  更新数据

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
			view.Render(w, view.D{
				"Article": _article,
				"Errors": errors,
			}, "articles.edit", "articles._form_field")
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

