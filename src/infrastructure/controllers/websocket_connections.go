package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/robertokbr/blinkchat/src/domain/enums"
	"github.com/robertokbr/blinkchat/src/domain/models"
	controller_errors "github.com/robertokbr/blinkchat/src/infrastructure/controllers/errors"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/infrastructure/websocket"
	"github.com/robertokbr/blinkchat/src/usecases"
)

type WebsocketConnections struct {
	Pool               *models.Pool
	UsersRepository    *repositories.Users
	SessionsRepository *repositories.Sessions
}

func (wsc *WebsocketConnections) Create(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")
	session, err := wsc.SessionsRepository.FindByToken(token)

	if err != nil {
		// handle error
	}

	if session == nil {
		// return not found error
	}

	user, err := wsc.UsersRepository.FindByID(session.UserID)

	if err != nil {
		// handle error
	}

	if user == nil {
		// return not found error
	}

	connection, err := websocket.Upgrade(r, w)

	if err != nil {
		controller_errors.WebsocketConnectionError(w, r, err)
		return
	}

	client := &models.Client{
		User:  user,
		Conn:  connection,
		State: enums.NOT_IN_A_MATCH,
	}

	registerClientUsecase := usecases.NewRegisterClient(wsc.Pool, client)
	unregisterClientUsecase := usecases.NewUnregisterClient(wsc.Pool, client)
	readClientMessagesUsecase := usecases.NewReadClientMessages(wsc.Pool, client, unregisterClientUsecase)

	registerClientUsecase.Execute()
	readClientMessagesUsecase.Execute()
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
