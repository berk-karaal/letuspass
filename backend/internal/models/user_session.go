package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model
	Token     string `gorm:"unique;index"`
	UserID    uint
	ExpiresAt time.Time `gorm:"not null"`
	UserAgent string
}
