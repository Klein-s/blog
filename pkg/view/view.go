package view

import (
	logger2 "goblog/pkg/logger"
	route "goblog/pkg/routes"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

/**
 	Render 渲染视图
 */
func Render(w io.Writer, name string, data interface{})  {
	//设置模板相对路径
	viewDir := "resources/views/"

	//将 articles.show 转为 articles/show
	name = strings.Replace(name, ".", "/", -1)
	//获取所有布局模板 slice
	files, err := filepath.Glob(viewDir+"layouts/*.gohtml")

	logger2.LogError(err)

	//将 slice加入 目标文件
	newFiles := append(files, viewDir+name+".gohtml")

	//解析所有模板文件
	tmpl, err := template.New(name + ".gohtml").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).
		ParseFiles(newFiles...)
	logger2.LogError(err)

	tmpl.ExecuteTemplate(w, "app", data)


}
