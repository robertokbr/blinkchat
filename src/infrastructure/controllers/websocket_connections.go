package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
	controller_errors "github.com/robertokbr/blinkchat/src/infrastructure/controllers/errors"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/infrastructure/pkg/websocket"
)

type WebsocketConnections struct {
	Pool            *websocket.Pool
	UsersRepository *repositories.Users
}

func (wsc *WebsocketConnections) Create(w http.ResponseWriter, r *http.Request) {
	connection, err := websocket.Upgrade(r, w)

	if err != nil {
		controller_errors.WebsocketConnectionError(w, r, err)
		return
	}

	query := r.URL.Query()

	user_id := query.Get("user_id")

	user, err := wsc.UsersRepository.FindByID(user_id)

	if err != nil {
		// return not found error
	}

	client := &websocket.Client{
		User:  user,
		Conn:  connection,
		Pool:  wsc.Pool,
		State: enums.NOT_IN_A_MATCH,
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
