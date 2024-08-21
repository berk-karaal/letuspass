package models

import "gorm.io/gorm"

const (
	VaultPermissionManageVault string = "manage_vault"
	VaultPermissionDeleteVault string = "delete_vault"
	VaultPermissionManageItems string = "manage_items"
	VaultPermissionRead        string = "read_vault"
)

type VaultPermission struct {
	gorm.Model
	VaultID    uint
	UserID     uint
	Permission string
}
