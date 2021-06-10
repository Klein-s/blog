package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "<h1>hello this is goblog</h1>")
}

func main()  {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":8005", nil)
}