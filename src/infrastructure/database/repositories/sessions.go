package repositories

import (
	"github.com/robertokbr/blinkchat/src/domain/models"
	"gorm.io/gorm"
)

type Sessions struct {
	db *gorm.DB
}

func NewSessions(db *gorm.DB) *Sessions {
	return &Sessions{db: db}
}

func (s *Sessions) Save(session *models.Session) error {
	err := s.db.Save(session).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *Sessions) FindByToken(token string) (*models.Session, error) {
	session := &models.Session{}

	err := s.db.Where("token = ?", token).First(session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}
