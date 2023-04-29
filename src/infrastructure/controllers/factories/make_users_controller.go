package controller_factories

import (
	"github.com/robertokbr/blinkchat/src/infrastructure/controllers"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
)

func MakeUsersController() *controllers.Users {
	usersRepo := repositories.NewUsers(database.Connection)
	sessionsRepo := repositories.NewSessions(database.Connection)
	usersController := controllers.NewUsers(usersRepo, sessionsRepo)

	return usersController
}
