package usecases

import (
	"errors"
	"log"

	"github.com/robertokbr/blinkchat/domain/interfaces"
	"github.com/robertokbr/blinkchat/domain/models"
)

type UpdateUser struct {
	UsersRepository interfaces.UsersRepository
}

func NewUpdateUser(usersRepository interfaces.UsersRepository) *UpdateUser {
	return &UpdateUser{UsersRepository: usersRepository}
}

func (uc *UpdateUser) Execute(user *models.User) error {
	user, err := uc.UsersRepository.FindByEmail(user.Email)

	if err != nil {
		log.Printf("error finding user: %v", err)
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	err = uc.UsersRepository.Save(user)

	if err != nil {
		log.Printf("error saving user: %v", err)
		return err
	}

	return nil
}
