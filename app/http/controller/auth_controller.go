package controller

import (
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct {

}

/**
	注册页面
 */
func (*AuthController) Register (w http.ResponseWriter, r *http.Request)  {
	view.RenderSimple(w, view.D{}, "auth.register")
}

/**
	注册
 */
func (*AuthController) DoRegister (w http.ResponseWriter, r *http.Request)  {

}