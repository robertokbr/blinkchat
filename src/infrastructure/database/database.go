package database

import (
	"github.com/robertokbr/blinkchat/src/domain/logger"
	"github.com/robertokbr/blinkchat/src/domain/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Dns           string
	AutoMigrateDb bool
	Debug         bool
}

func NewDatabase() *Database {
	return &Database{
		Dns:           ":memory:",
		AutoMigrateDb: true,
		Debug:         true,
	}
}

var Connection *gorm.DB

func (db *Database) Connect() (*gorm.DB, error) {
	if Connection != nil {
		return Connection, nil
	}

	logger.Infof("connecting to database: %s", db.Dns)

	var err error

	Connection, err = gorm.Open(sqlite.Open(db.Dns))

	if err != nil {
		return &gorm.DB{}, err
	}

	if db.AutoMigrateDb {
		err := Connection.AutoMigrate(
			&models.User{},
			&models.Session{},
		)

		if err != nil {
			logger.Infof("error automigrating database: %v", err)
			return &gorm.DB{}, err
		}
	}

	return Connection, nil
}
