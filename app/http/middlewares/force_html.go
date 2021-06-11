package middlewares

import "net/http"

func ForceHTML(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//设置请求标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		//继续处理请求
		next.ServeHTTP(w, r)
	})
}
