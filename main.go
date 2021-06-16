package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	c "goblog/pkg/config"
	"net/http"
)


var router *mux.Router

func init()  {
	config.Initialize()
}

func main() {

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()


	http.ListenAndServe(":" + c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
}