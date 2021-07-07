package view

import (
	"goblog/app/model/category"
	"goblog/app/model/user"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	logger2 "goblog/pkg/logger"
	route "goblog/pkg/routes"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

type D map[string]interface{}
/**
 	Render 渲染通用视图
 */
func Render(w io.Writer,  data D, tplFiles ...string)  {
	renderTemplate(w, "app", data, tplFiles...)
}

/**
	RenderSimple 渲染简单视图
*/
func RenderSimple(w io.Writer,  data D, tplFiles ...string)  {
	renderTemplate(w, "simple", data, tplFiles...)
}

/**
	RenderTemplate 渲染视图
*/
func renderTemplate(w io.Writer, name string,  data D, tplFiles ...string)  {

	data["isLogined"] = auth.Check()
	data["flash"] = flash.All()
	data["loginuser"] = auth.User()
	data["Categories"], _ = category.All()
	data["Users"], _ = user.Limit3()
	data["IsAdmin"] = auth.User().HasRole("admin")

	//设置模板相对路径
	viewDir := "resources/views/"

	//遍历 tplFiles,设置正确路径
	for i, f := range tplFiles{
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	//获取所有布局模板 slice
	layoutFiles, err := filepath.Glob(viewDir+"layouts/*.gohtml")

	logger2.LogError(err)

	//将 slice加入 目标文件
	allFiles := append(layoutFiles, tplFiles...)

	//解析所有模板文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Html": Html,
		}).
		ParseFiles(allFiles...)
	logger2.LogError(err)

	tmpl.ExecuteTemplate(w, name, data)
}

func Html(content string) interface{} {
	return template.HTML(content)
}