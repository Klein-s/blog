package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/pkg/database"
	logger2 "goblog/pkg/logger"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

var router *mux.Router

var db *sql.DB


type Article struct {
	Title, Body string
	ID int64
}




func (a Article) Delete() (rowsAffected int64, err error)  {
	
	rs, err := db.Exec("delete from articles where id = " + strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}
	//删除成功 跳转到文章详情页面
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}



func getArticleByID(id string) (Article, error)  {
	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request)  {

	//获取url参数

	id := getRouteVariable("id", r)

	//读取对应的文章数据
	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
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
		updateUrl, _ := router.Get("articles.update").URL("id", id)
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
	更新文章
 */
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request)  {

	//获取url参数
	id := getRouteVariable("id", r)

	//获取对应的文章数据
	_, err := getArticleByID(id)

	//处理错误
	if err != nil {
		if err == sql.ErrNoRows {
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

			query := "update articles set title = ? , body = ? where id = ?"
			rs, err := db.Exec(query, title, body, id)
			if err != nil {
				logger2.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			//更新成功，跳转到文章详情页面

			if n, _ := rs.RowsAffected(); n > 0 {
				showUrl, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showUrl.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您未做任何修改")
			}
		} else {
			 //表单验证未通过

			updateUrl, _ := router.Get("articles.update").URL("id", id)
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
	删除文章
 */

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request)  {
	
	//获取url参数
	id := getRouteVariable("id", r)
	
	//获取文章数据
	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
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
		rowsAffected, err := article.Delete()
		
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
				indexUrl, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexUrl.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
		
	}
}



/**
	文章列表
 */
func articlesIndexHandler(w http.ResponseWriter, r *http.Request)  {

	//执行查询语句，返回结果集
	rows, err := db.Query("select * from articles")
	logger2.LogError(err)
	defer rows.Close()

	var articles []Article
	//循环读取结果
	for rows.Next() {
		var article Article
		//扫描每行的结果 并赋值到 article 对象中
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		logger2.LogError(err)
		//将 article 追加到 articles 数组中
		articles = append(articles, article)
	}

	//检查遍历时是否发生错误
	err = rows.Err()
	logger2.LogError(err)

	//加载模板
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	logger2.LogError(err)

	//渲染模板，将所有文章数据传输进去
	tmpl.Execute(w, articles)
}

type ArticlesFormData struct {
	Title, Body string
	URL *url.URL
	Errors map[string]string
}

/**
	创建文章
 */
func articlesStoreHandler(w http.ResponseWriter, r *http.Request)  {

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)


	//检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			logger2.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}

	} else {

		storeURl, _ := router.Get("articles.store").URL()

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
func saveArticleToDB(title string, body string) (int64, error)  {

	//变量初始化
	var (
		id int64
		err error
		rs sql.Result
		stmt *sql.Stmt
	)

	//获取一个 prepare 声明语句
	stmt, err = db.Prepare("insert into articles(title, body) value(?,?)")
	//错误检测
	if err != nil {
		return 0, err
	}

	//函数运行后关闭， 防止占用sql连接
	defer stmt.Close()

	//执行请求， 传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	// 插入成功，返回自增ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, err
}

func forceHTMLMiddleware(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//设置请求标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		//继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request)  {

	storeUrl, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title: "",
		Body: "",
		URL: storeUrl,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
}

//获取url 参数
func getRouteVariable(parameterName string, r *http.Request) string  {
	vars := mux.Vars(r)
	return vars[parameterName]
}


func main() {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()



	router.HandleFunc("/articles", articlesStoreHandler).
		Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).
		Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).
		Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).
		Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).
		Methods("POST").Name("articles.delete")


	//中间件，强制内容为HTMl
	router.Use(forceHTMLMiddleware)
	//通过命名路由获取URL示例
	//homeURl, _ := router.Get("home").URL()
	//fmt.Println("homeURL:", homeURl)
	//articleURL, _ := router.Get("articles.show").
	//	URL("id", "23")
	//fmt.Println("articleURL", articleURL)

	http.ListenAndServe(":8005", removeTrailingSlash(router))
}