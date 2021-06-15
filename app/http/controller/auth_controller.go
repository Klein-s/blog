package controller

import (
	"fmt"
	"goblog/app/model/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
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

	//初始化数据
	_user := user.User{
		Name: r.PostFormValue("name"),
		Email: r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		//data, _ := json.MarshalIndent(errs, "", " ")
		//fmt.Fprint(w, string(data))
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User": _user,
		}, "auth.register")
	} else {
		_user.Create()

		if _user.ID > 0 {
			fmt.Fprint(w, "插入成功 ID 为" + _user.GetStringID())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建失败")
		}
	}

}

func (*AuthController) Login(w http.ResponseWriter, r *http.Request)  {
	//fmt.Fprint(w, session.Get("uid"))
	view.RenderSimple(w, view.D{}, "auth.login")
}

func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request)  {

	//获取表单数据
	email :=  r.PostFormValue("email")
	password :=  r.PostFormValue("password")

	//尝试登录
	if err := auth.Attempt(email, password); err == nil {
		//登录成功
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		//失败 显示错误提示
		view.RenderSimple(w, view.D{
			"Error": err.Error(),
			"Email": email,
			"Password": password,
		}, "auth.login")
	}
}

func (*AuthController) Logout(w http.ResponseWriter, r *http.Request)  {
	auth.Logout()
	http.Redirect(w, r, "/", http.StatusFound)
}