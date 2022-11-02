package usecase_tests

import (
	"testing"

	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/database"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/usecases"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db := database.NewDatabase()

	connection, err := db.Connect()

	require.Nil(t, err)

	usersRepository := repositories.NewUsersRepository(connection)

	createUser := usecases.NewCreateUser(usersRepository)

	createUserDTO := dtos.CreateConnection{
		Email:    "email@email.com",
		Name:     "name",
		ImageURL: "image_url",
	}

	createdUser, err := createUser.Execute(createUserDTO)

	require.Nil(t, err)

	require.Equal(t, createUserDTO.Email, createdUser.Email)
}

func TestReturnAlreadyCreatedUser(t *testing.T) {
	db := database.NewDatabase()
	connection, err := db.Connect()

	require.Nil(t, err)

	usersRepository := repositories.NewUsersRepository(connection)
	createUser := usecases.NewCreateUser(usersRepository)

	createUserDTO := dtos.CreateConnection{
		Email:    "email@email.com",
		Name:     "name",
		ImageURL: "image_url",
	}

	user := models.NewUser(createUserDTO)

	usersRepository.Save(user)

	createdUser, err := createUser.Execute(createUserDTO)

	require.Nil(t, err)

	require.Equal(t, user.Email, createdUser.Email)
}
