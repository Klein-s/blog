package controller

import (
	"fmt"
	"goblog/pkg/flash"
	logger2 "goblog/pkg/logger"
	"gorm.io/gorm"
	"net/http"
)

type BaseController struct {

}

// 处理sql错误并返回
func (bc BaseController) ResponseForSQLError(w http.ResponseWriter, err error)  {
	if err == gorm.ErrRecordNotFound {
		// 3.1 数据未找到
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 文章未找到")
	} else {
		// 3.2 数据库错误
		logger2.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	}
}

func (bc BaseController) ResponseForUnauthorized(w http.ResponseWriter, r *http.Request)  {
	flash.Warning("未授权操作！")
	http.Redirect(w, r, "/", http.StatusFound)
}