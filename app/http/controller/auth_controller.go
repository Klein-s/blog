package controller

import (
	"fmt"
	"goblog/app/model/user"
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

	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	//表单验证
	//验证通过  入库  跳转首页

	_user := user.User{
		Name: name,
		Email: email,
		Password: password,
	}
	_user.Create()

	if _user.ID > 0 {
		fmt.Fprint(w, "插入成功 ID 为" + _user.GetStringID())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "创建失败")
	}
	//表单不通过  重新显示页面
}