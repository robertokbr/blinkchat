package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/robertokbr/blinkchat/src/domain/errors"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/infrastructure/websocket"
	"github.com/robertokbr/blinkchat/src/usecases"
)

type WebsocketClients struct {
	pool               *models.Pool
	jobs               chan models.Message
	usersRepository    *repositories.Users
	sessionsRepository *repositories.Sessions
}

func NewWebsocketClients(
	pool *models.Pool,
	jobs chan models.Message,
	usersRepository *repositories.Users,
	sessionsRepository *repositories.Sessions,
) *WebsocketClients {
	return &WebsocketClients{
		jobs:               jobs,
		pool:               pool,
		usersRepository:    usersRepository,
		sessionsRepository: sessionsRepository,
	}
}

func (wsc *WebsocketClients) Create(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	token := query.Get("token")
	findSessionUserUsecase := usecases.NewFindSessionsUser(wsc.sessionsRepository, wsc.usersRepository)
	user, err := findSessionUserUsecase.Execute(token)

	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	connection, err := websocket.Upgrade(r, w)

	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	client := models.NewClient(user, connection)
	registerClientUsecase := usecases.NewRegisterClient(client, wsc.pool)
	unregisterClientUsecase := usecases.NewUnregisterClient(client, wsc.pool)
	readClientMessagesUsecase := usecases.NewReadClientMessages(client, wsc.jobs)

	defer unregisterClientUsecase.Execute()
	defer readClientMessagesUsecase.Stop()

	registerClientUsecase.Execute()
	readClientMessagesUsecase.Execute()
}

func (wsc *WebsocketClients) FindAll(w http.ResponseWriter, r *http.Request) {
	connections := make([]*models.User, 0)

	for _, client := range wsc.pool.Clients {
		connections = append(connections, client.User)
	}

	serialized, err := json.Marshal(connections)

	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Write(serialized)
}
