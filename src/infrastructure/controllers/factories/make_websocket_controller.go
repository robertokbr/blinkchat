package controller_factories

import (
	"github.com/robertokbr/blinkchat/src/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/infrastructure/pkg/websocket"
)

func MakeWebsocketConnectionsController() *controllers.WebsocketConnections {
	usersRepo := repositories.NewUsersRepository(database.Connection)

	var pool = websocket.NewPool()

	var websocketConnectionsController = controllers.WebsocketConnections{
		Pool:            pool,
		UsersRepository: usersRepo,
	}

	return &websocketConnectionsController
}
