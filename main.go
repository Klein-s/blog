package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/pkg/database"
	logger2 "goblog/pkg/logger"
	"net/http"
	"strconv"
	"strings"
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





	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).
		Methods("POST").Name("articles.delete")


	//中间件，强制内容为HTMl
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":8005", removeTrailingSlash(router))
}