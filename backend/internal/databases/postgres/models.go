package postgres

import "github.com/berk-karaal/letuspass/backend/internal/models"

func GetModels() []interface{} {
	return []interface{}{
		&models.User{}, &models.UserSession{}, &models.Vault{}, &models.VaultPermission{},
		&models.VaultItem{}, &models.VaultKey{}, &models.VaultAuditLog{},
	}
}
