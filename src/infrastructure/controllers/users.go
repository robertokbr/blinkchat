package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/robertokbr/blinkchat/src/domain/dtos"
	"github.com/robertokbr/blinkchat/src/domain/errors"
	"github.com/robertokbr/blinkchat/src/infrastructure/database/repositories"
	"github.com/robertokbr/blinkchat/src/usecases"
)

type Users struct {
	usersRepository    *repositories.Users
	sessionsRepository *repositories.Sessions
}

func NewUsers(usersRepository *repositories.Users, sessionsRepository *repositories.Sessions) *Users {
	return &Users{
		usersRepository:    usersRepository,
		sessionsRepository: sessionsRepository,
	}
}

func (wsc *Users) Create(w http.ResponseWriter, r *http.Request) {
	createUserDTO := dtos.CreateUser{}

	err := json.NewDecoder(r.Body).Decode(&createUserDTO)

	defer r.Body.Close()

	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	createUserUsecase := usecases.NewCreateUser(wsc.usersRepository, wsc.sessionsRepository)
	session, err := createUserUsecase.Execute(&createUserDTO)

	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	sessionEncoded, err := json.Marshal(session)

	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(sessionEncoded)
}
