package usecases

import (
	"github.com/robertokbr/blinkchat/src/domain/errors"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
)

type FindSessionsUser struct {
	sessionsRepository *repositories.Sessions
	usersRepository    *repositories.Users
}

func NewFindSessionsUser(
	sessionsRepository *repositories.Sessions,
	usersRepository *repositories.Users,
) *FindSessionsUser {
	return &FindSessionsUser{
		sessionsRepository: sessionsRepository,
		usersRepository:    usersRepository,
	}
}

func (uc *FindSessionsUser) Execute(token string) (*models.User, error) {
	session, err := uc.sessionsRepository.FindByToken(token)

	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.Unauthorized("invalid token")
	}

	if session.IsExpired() {
		return nil, errors.Unauthorized("session expired")
	}

	user, err := uc.usersRepository.FindByID(session.UserID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.NotFound("user not found to this session")
	}

	return user, nil
}
