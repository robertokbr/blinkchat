package middlewares

import "net/http"

func Cors(r *http.Request) bool {
	return true
}
