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
