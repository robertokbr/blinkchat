package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/robertokbr/blinkchat/infrastructure/middlewares"
)

func Upgrade(r *http.Request, w http.ResponseWriter) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     middlewares.Cors,
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return &websocket.Conn{}, err
	}

	return conn, nil
}
