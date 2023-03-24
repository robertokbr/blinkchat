package usecase_factories

import (
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/usecases"
)

func MakeCreateUserUsecase() *usecases.CreateUser {
	usersRepo := repositories.NewUsersRepository(database.Connection)

	var createUserUsecase = usecases.NewCreateUser(usersRepo)

	return createUserUsecase
}
