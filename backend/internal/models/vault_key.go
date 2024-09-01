package models

import "gorm.io/gorm"

type VaultKey struct {
	gorm.Model
	VaultID           uint
	KeyOwnerUserID    uint
	InviterUserID     uint
	EncryptionIV      string
	EncryptedVaultKey string

	KeyOwnerUser User `gorm:"foreignKey:KeyOwnerUserID"`
	InviterUser  User `gorm:"foreignKey:InviterUserID"`
}
