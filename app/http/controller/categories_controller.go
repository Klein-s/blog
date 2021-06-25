package controller

import (
	"fmt"
	"goblog/app/model/category"
	"goblog/app/requests"
	"goblog/pkg/view"
	"net/http"
)

type CategoriesController struct {
	BaseController
}

//文章分类创建页面
func (* CategoriesController) Create(w http.ResponseWriter, r *http.Request)  {
	view.Render(w, view.D{}, "categories.create")
}

//文章分类数据保存
func (* CategoriesController) Store(w http.ResponseWriter, r *http.Request)  {
//fmt.Fprint(w,"123123")
	//初始化数据
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}

	//表单验证
	errors := requests.ValidateCategoryForm(_category)

	//检查是否有错误
	if len(errors) == 0 {

		_category.Create()
		if _category.ID > 0 {
			//fmt.Fprint(w, "插入成功，ID 为"+_article.GetStringID())
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建失败，请联系管理员")
		}

	} else {

		view.Render(w, view.D{
			"Category": _category,
			"Errors":  errors,
		}, "categories.create")
	}

}
