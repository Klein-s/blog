package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

var router = mux.NewRouter()


func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>hello 欢迎来到 goblog</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"994097656@qq.com\">994097656@qq.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID: ", id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "访问文章列表")
}

type ArticlesFormData struct {
	Title, Body string
	URL *url.URL
	Errors map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request)  {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	//验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度为 3-40"
	}

	//验证内容
	if body == "" {
		errors["body"] = "内容不可为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度必须大于或等于10个字节"
	}

	//检查是否有错误
	if len(errors) == 0 {
		fmt.Fprint(w, "通过验证！<br>")
		fmt.Fprintf(w, "title 的值为: %v <br>", title)
		fmt.Fprintf(w, "title 的长度为: %v <br>", len(title))
		fmt.Fprintf(w, "body 的值为: %v <br>", body)
		fmt.Fprintf(w, "body 的长度为: %v <br>", len(body))
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

func main() {


	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}",
		articlesShowHandler).
		Methods("GET").
		Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).
		Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).
		Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

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