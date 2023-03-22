package factories

import (
	"github.com/robertokbr/blinkchat/src/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/infrastructure/pkg/websocket"
	"github.com/robertokbr/blinkchat/src/usecases"
)

func MakeWebsocketConnectionsController() *controllers.WebsocketConnections {
	usersRepo := repositories.NewUsersRepository(database.Connection)

	createUserUC := usecases.NewCreateUser(usersRepo)

	var pool = websocket.NewPool()

	var websocketConnectionsController = controllers.WebsocketConnections{
		Pool:              pool,
		CreateUserUsecase: createUserUC,
	}

	return &websocketConnectionsController
}
