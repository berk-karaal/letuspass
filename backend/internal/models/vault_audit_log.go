package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLogAction string

type VaultAuditLog struct {
	gorm.Model
	VaultID     uint
	VaultItemID uint `gorm:"default:null"`
	UserID      uint
	ActionCode  AuditLogAction
	ActionData  datatypes.JSONMap

	Vault     Vault     `gorm:"foreignKey:VaultID"`
	VaultItem VaultItem `gorm:"foreignKey:VaultItemID"`
	User      User      `gorm:"foreignKey:UserID"`
}

const (
	AuditLogActionVaultCreate     AuditLogAction = "vault_create"
	AuditLogActionVaultRename     AuditLogAction = "vault_rename"
	AuditLogActionVaultDelete     AuditLogAction = "vault_delete"
	AuditLogActionVaultAddUser    AuditLogAction = "vault_add_user"
	AuditLogActionVaultRemoveUser AuditLogAction = "vault_remove_user"
	AuditLogActionVaultUserLeft   AuditLogAction = "vault_user_left"
	AuditLogActionVaultItemCreate AuditLogAction = "vault_item_create"
	AuditLogActionVaultItemUpdate AuditLogAction = "vault_item_update"
	AuditLogActionVaultItemDelete AuditLogAction = "vault_item_delete"
)

func AuditLogDataVaultCreate(name string) map[string]any {
	return map[string]any{
		"name": name,
	}
}

func AuditLogDataVaultRename(oldName, newName string) map[string]any {
	return map[string]any{
		"old_name": oldName,
		"new_name": newName,
	}
}

func AuditLogDataVaultDelete() map[string]any {
	return map[string]any{}
}
