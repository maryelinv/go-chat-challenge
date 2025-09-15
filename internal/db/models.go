package db

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash []byte `gorm:"not null"`
	CreatedAt    time.Time
}

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	Room      string `gorm:"index;not null"`
	UserID    *uint
	Username  string    `gorm:"not null"`
	Text      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"index"`
}

type Room struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	Creator   string `gorm:"index"`
	IsPrivate bool   `gorm:"default:false;not null"`
	CreatedAt time.Time
}
