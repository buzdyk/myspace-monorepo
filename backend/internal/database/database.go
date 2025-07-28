package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Tracker{},
		&Project{},
		&Track{},
		&Setting{},
	)
}