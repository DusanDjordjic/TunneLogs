package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB = nil

func Connect() error {
	db, err := gorm.Open(sqlite.Open("tunnelogs.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	DB.AutoMigrate()
	return nil
}
