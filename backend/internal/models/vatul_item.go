package models

import "gorm.io/gorm"

type VaultItem struct {
	gorm.Model
	VaultID           uint
	Title             string
	EncryptedUsername string
	EncryptedPassword string
	EncryptedNote     string
}
