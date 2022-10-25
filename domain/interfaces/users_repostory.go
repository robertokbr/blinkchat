package interfaces

import (
	"github.com/robertokbr/blinkchat/domain/models"
)

type UsersRepository interface {
	FindByEmail(email string) (*models.User, error)
	Save(user *models.User) error
}
