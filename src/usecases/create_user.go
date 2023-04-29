package usecases

import (
	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/errors"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
)

type CreateUser struct {
	userRepository     *repositories.Users
	sessionsRepository *repositories.Sessions
}

func NewCreateUser(userRepository *repositories.Users, sessionsRepository *repositories.Sessions) *CreateUser {
	return &CreateUser{
		userRepository:     userRepository,
		sessionsRepository: sessionsRepository,
	}
}

func (c *CreateUser) Execute(data *dtos.CreateUser) (*models.Session, error) {
	user, err := c.userRepository.FindByEmail(data.Email)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, errors.Conflict("user already exists with this email")
	}

	user = models.NewUser(data)

	err = c.userRepository.Save(user)

	if err != nil {
		return nil, err
	}

	session := models.NewSession(user.ID)

	err = c.sessionsRepository.Save(session)

	if err != nil {
		return nil, err
	}

	return session, nil
}
