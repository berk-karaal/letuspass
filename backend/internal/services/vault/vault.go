package vault

import (
	"github.com/berk-karaal/letuspass/backend/internal/models"
	"gorm.io/gorm"
)

// CheckUserHasVaultPermission returns true if given user has given permission on given vault, if not returns false.
func CheckUserHasVaultPermission(db *gorm.DB, userId, vaultId int, permission string) (hasPermission bool, err error) {
	err = db.Model(&models.VaultPermission{}).Select("count(*) > 0").
		Where("vault_id = ? AND user_id = ? AND permission = ?", vaultId, userId, permission).
		Scan(&hasPermission).Error
	if err != nil {
		return false, err
	}
	return hasPermission, nil
}
