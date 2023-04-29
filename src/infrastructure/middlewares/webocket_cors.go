package middlewares

import "net/http"

func WebsocketCors(r *http.Request) bool {
	return true
}
