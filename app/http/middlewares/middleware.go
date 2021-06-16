package middlewares

import "net/http"

type HttpHandlerFunc func(w http.ResponseWriter, r *http.Request)
