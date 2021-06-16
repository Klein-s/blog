package middlewares

import (
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"net/http"
)

// 未登录用户才能访问
func Guest(next http.HandlerFunc) HttpHandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		if auth.Check() {
			flash.Warning("登录用户无法访问此页面")
			http.Redirect(w, r ,"/", http.StatusFound)
			return
		}

		next(w, r)
	}
}
