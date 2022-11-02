package repositories

import (
	"github.com/robertokbr/blinkchat/src/domain/models"
	"gorm.io/gorm"
)

type Users struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *Users {
	return &Users{db: db}
}

func (u *Users) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := u.db.Where("email = ?", email).First(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (u *Users) Save(user *models.User) error {
	err := u.db.Save(user).Error

	if err != nil {
		return err
	}

	return nil
}
