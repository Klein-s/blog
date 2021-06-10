package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter()

var db *sql.DB

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

type Article struct {
	Title, Body string
	ID int64
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request)  {

	//获取url参数
	vars := mux.Vars(r)
	id := vars["id"]

	//读取对应文章数据
	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	//处理错误
	if err != nil {
		if err == sql.ErrNoRows {
			//数据为找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			//数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		//读取成功
		tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w, article)
	}
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
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
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

//数据库驱动初始化
func initDB()  {

	var err error
	config := mysql.Config{
		User: "rdbuser",
		Passwd: "shop2db123#",
		Addr: "101.37.150.149:3306",
		Net: "tcp",
		DBName: "blog",
		AllowNativePasswords: true,
	}

	//准备数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	//设置最大连接数
	db.SetMaxOpenConns(25)
	//设置最大空闲连接数
	db.SetMaxIdleConns(25)
	//设置每个连接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	//尝试连接 失败报错
	err = db.Ping()
	checkError(err)
}

func checkError(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func main() {
	initDB()
	createTables()

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