package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"
)


var router *mux.Router


func main() {

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()


	http.ListenAndServe(":8005", middlewares.RemoveTrailingSlash(router))
}