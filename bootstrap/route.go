package bootstrap

import (
	"github.com/gorilla/mux"
	route "goblog/pkg/routes"
	"goblog/routes"
)

func SetupRoute() *mux.Router  {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	route.SetRoute(router)
	return router
}