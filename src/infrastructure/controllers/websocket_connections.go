package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/pkg/websocket"
	"github.com/robertokbr/blinkchat/src/usecases"
)

type WebsocketConnections struct {
	Pool              *websocket.Pool
	CreateUserUsecase *usecases.CreateUser
}

func (wsc *WebsocketConnections) Create(w http.ResponseWriter, r *http.Request) {
	connection, err := websocket.Upgrade(r, w)

	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not upgrade connection"))
		return
	}

	var createConnectionDTO dtos.CreateConnection

	query := r.URL.Query()

	createUserString := query.Get("user")

	json.Unmarshal([]byte(createUserString), &createConnectionDTO)

	log.Printf("Connecting user: %v", createConnectionDTO)

	user, err := wsc.CreateUserUsecase.Execute(createConnectionDTO)

	client := &websocket.Client{
		User: user,
		Conn: connection,
		Pool: wsc.Pool,
	}

	wsc.Pool.Register(client)

	client.Read()
}

// Return all connected users
func (wsc *WebsocketConnections) FindAll(w http.ResponseWriter, r *http.Request) {
	connections := make([]*models.User, 0)

	for _, client := range wsc.Pool.Clients {
		connections = append(connections, client.User)
	}

	serialized, err := json.Marshal(connections)

	if err != nil {
		log.Printf("error serializing connections: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not serialize connections"))
		return
	}

	w.Write(serialized)
}
