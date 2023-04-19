package controller_factories

import (
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
)

func MakeWebsocketConnectionsController(pool *models.Pool) *controllers.WebsocketConnections {
	usersRepo := repositories.NewUsersRepository(database.Connection)
	sessionsRepo := repositories.NewSessionsRepository(database.Connection)
	websocketConnectionsController := controllers.WebsocketConnections{
		Pool:               pool,
		UsersRepository:    usersRepo,
		SessionsRepository: sessionsRepo,
	}

	return &websocketConnectionsController
}
