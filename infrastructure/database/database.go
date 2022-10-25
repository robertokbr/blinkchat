package database

import (
	"log"

	"github.com/robertokbr/blinkchat/domain/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Db            *gorm.DB
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

func (db *Database) Connect() (*gorm.DB, error) {
	var err error

	db.Db, err = gorm.Open(sqlite.Open(db.Dns))

	if err != nil {
		return &gorm.DB{}, err
	}

	if db.AutoMigrateDb {
		err := db.Db.AutoMigrate(
			&models.User{},
		)

		if err != nil {
			log.Printf("error automigrating database: %v", err)
			return &gorm.DB{}, err
		}
	}

	return db.Db, nil
}
