package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email             string `gorm:"unique;index"`
	Password          string
	Name              string
	KeyDerivationSalt string
	PublicKey         string
	IsActive          bool `gorm:"default:true"`

	UserSessions     []UserSession
	VaultPermissions []VaultPermission
}
