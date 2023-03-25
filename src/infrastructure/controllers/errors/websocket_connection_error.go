package controller_errors

import (
	"net/http"

	"github.com/robertokbr/blinkchat/src/domain/logger"
)

func WebsocketConnectionError(w http.ResponseWriter, r *http.Request, err error) {
	logger.Errorf("error upgrading connection: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Could not upgrade connection"))
}
