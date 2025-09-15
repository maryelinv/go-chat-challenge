package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path), &gorm.Config{})
}

func Migrate(g *gorm.DB) error {
	return g.AutoMigrate(&User{}, &Message{}, &Room{})
}
