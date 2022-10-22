package interfaces

import (
	"github.com/robertokbr/blinkchat/domain/dtos"
	"github.com/robertokbr/blinkchat/domain/models"
)

type UsersRepository interface {
	Create(data dtos.CreateConnection) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Save(user *models.User) error
}
