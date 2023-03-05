package usecases

import (
	"log"

	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/interfaces"
	"github.com/robertokbr/blinkchat/src/domain/models"
)

type CreateUser struct {
	UsersRepository interfaces.UsersRepository
}

func NewCreateUser(usersRepository interfaces.UsersRepository) *CreateUser {
	return &CreateUser{UsersRepository: usersRepository}
}

// TODO: should have a route of prelogin to check if the user is the true owner of the email
// This route should send an email to the user with a token to confirm the email
// This token should be used in the CreateUser route
func (uc *CreateUser) Execute(data dtos.CreateConnection) (*models.User, error) {
	user, err := uc.UsersRepository.FindByEmail(data.Email)

	if err != nil {
		log.Printf("error finding user: %v", err)
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user = models.NewUser(data)

	err = uc.UsersRepository.Save(user)

	if err != nil {
		log.Printf("error saving user: %v", err)
		return nil, err
	}

	return user, nil
}
