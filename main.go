package main

import (
	"fmt"
	"net/http"
)

func defaultHandFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
			"<p>如有疑惑，请联系我们。</p>")
	}
}

func aboutHandFunc(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"994097656@qq.com\">994097656@qq.com</a>")
}
func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", defaultHandFunc)
	router.HandleFunc("/about", aboutHandFunc)
	http.ListenAndServe(":8005", router)
}