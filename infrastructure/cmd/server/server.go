package main

import (
	"log"
	"net/http"

	"github.com/robertokbr/blinkchat/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/infrastructure/pkg/websocket"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	var pool = websocket.NewPool()

	var websocketConnectionsController = controllers.WebsocketConnections{Pool: pool}

	for i := 0; i < 10; i++ {
		go pool.Start(i)
	}

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/ws", websocketConnectionsController.Create)
	http.HandleFunc("/connections", websocketConnectionsController.FindAll)

	log.Println("Server started on port 8080")

	http.ListenAndServe(":8080", nil)
}
