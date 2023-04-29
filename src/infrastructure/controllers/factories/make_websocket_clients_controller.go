package controller_factories

import (
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
)

func MakeWebsocketClientsController(pool *models.Pool, jobs chan models.Message) *controllers.WebsocketClients {
	usersRepo := repositories.NewUsers(database.Connection)
	sessionsRepo := repositories.NewSessions(database.Connection)
	websocketClientsController := controllers.NewWebsocketClients(pool, jobs, usersRepo, sessionsRepo)

	return websocketClientsController
}
