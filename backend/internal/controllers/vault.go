package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/berk-karaal/letuspass/backend/internal/common"
	"github.com/berk-karaal/letuspass/backend/internal/common/bodybinder"
	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/common/orderbyparam"
	"github.com/berk-karaal/letuspass/backend/internal/common/pagination"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/models"
	"github.com/berk-karaal/letuspass/backend/internal/schemas"
	vaultservice "github.com/berk-karaal/letuspass/backend/internal/services/vault"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// HandleVaultsCreate
//
//	@Summary	Create a new vault
//	@Tags		vaults
//	@Id			createVault
//	@Param		request	body	controllers.HandleVaultsCreate.VaultCreateRequest	true	"New vault data"
//	@Produce	json
//	@Success	201	{object}	controllers.HandleVaultsCreate.VaultCreateResponse
//	@Failure	401
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults [post]
func HandleVaultsCreate(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultCreateRequest struct {
		Name              string `json:"name" binding:"required"`
		EncryptionIV      string `json:"encryption_iv" binding:"required"`
		EncryptedVaultKey string `json:"encrypted_vault_key" binding:"required"`
	}

	type VaultCreateResponse struct {
		Id   uint   `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	return func(c *gin.Context) {
		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		var requestData VaultCreateRequest
		if ok = bodybinder.Bind(&requestData, c); !ok {
			return
		}

		vault := models.Vault{Name: requestData.Name}
		if err := db.Create(&vault).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating vault failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		vaultKey := models.VaultKey{
			VaultID:           vault.ID,
			KeyOwnerUserID:    user.ID,
			InviterUserID:     user.ID,
			EncryptionIV:      requestData.EncryptionIV,
			EncryptedVaultKey: requestData.EncryptedVaultKey,
		}
		if err := db.Create(&vaultKey).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating vault key failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		vaultPermissions := []models.VaultPermission{
			{VaultID: vault.ID, UserID: user.ID, Permission: models.VaultPermissionManageVault},
			{VaultID: vault.ID, UserID: user.ID, Permission: models.VaultPermissionDeleteVault},
			{VaultID: vault.ID, UserID: user.ID, Permission: models.VaultPermissionManageItems},
			{VaultID: vault.ID, UserID: user.ID, Permission: models.VaultPermissionRead},
		}
		if err := db.Create(&vaultPermissions).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating vault permissions failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		auditLog := models.VaultAuditLog{
			VaultID:     vault.ID,
			VaultItemID: 0,
			UserID:      user.ID,
			ActionCode:  models.AuditLogActionVaultCreate,
			ActionData:  models.AuditLogDataVaultCreate(vault.Name),
		}
		if err := db.Create(&auditLog).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Saving audit log failed.")
		}

		c.JSON(http.StatusCreated, VaultCreateResponse{Id: vault.ID, Name: vault.Name})
	}
}

// HandleVaultsList
//
//	@Summary	List vaults that user has read access to
//	@Tags		vaults
//	@Id			listVaults
//	@Produce	json
//	@Param		page		query		int		false	"Page number"			default(1)	minimum(1)
//	@Param		page_size	query		int		false	"Item count per page"	default(10)
//	@Param		ordering	query		string	false	"Ordering"				Enums(name, -name, created_at, -created_at)
//	@Success	200			{object}	pagination.StandardPaginationResponse[controllers.HandleVaultsList.VaultResponseItem]
//	@Failure	401
//	@Failure	500
//	@Router		/vaults [get]
func HandleVaultsList(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultResponseItem struct {
		Id        uint      `json:"id" binding:"required"`
		Name      string    `json:"name" binding:"required"`
		CreatedAt time.Time `json:"created_at" binding:"required"`
		UpdatedAt time.Time `json:"updated_at" binding:"required"`
	}

	return func(c *gin.Context) {
		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		ordering, err := orderbyparam.GenerateOrdering(c, map[string]string{
			"name":       "vaults.name",
			"created_at": "vaults.created_at",
		}, "name")
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Generating query ordering from params failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		var count int64
		err = db.Select("count(*)").
			Model(models.VaultPermission{}).
			Where("vault_permissions.user_id = ? AND vault_permissions.permission = ?",
				user.ID, models.VaultPermissionRead).
			Count(&count).Error

		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying user's vaults count failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		results := []VaultResponseItem{}
		err = db.Scopes(pagination.Paginate(c)).Select("vaults.id, vaults.name, vaults.created_at, vaults.updated_at").
			Model(models.VaultPermission{}).
			Joins("LEFT OUTER JOIN vaults ON vault_permissions.vault_id = vaults.id").
			Where("vault_permissions.user_id = ? AND vault_permissions.permission = ?", user.ID, models.VaultPermissionRead).
			Order(ordering).Scan(&results).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying user's vaults failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, pagination.StandardPaginationResponse[VaultResponseItem]{
			Results: results,
			Count:   int(count),
		})
	}
}

// HandleVaultsRetrieve
//
//	@Summary	Retrieve vault by id
//	@Tags		vaults
//	@Id			retrieveVault
//	@Produce	json
//	@Success	200	{object}	controllers.HandleVaultsCreate.VaultCreateResponse
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Forbidden	403
//	@Failure	500
//	@Router		/vaults/{id} [get]
//	@Param		id	path	int	true	"Vault id"
func HandleVaultsRetrieve(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {

	type VaultRetrieveResponse struct {
		Id        uint      `json:"id" binding:"required"`
		Name      string    `json:"name" binding:"required"`
		CreatedAt time.Time `json:"created_at" binding:"required"`
	}

	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		canRead, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionRead)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canRead {
			c.Status(http.StatusForbidden)
			return
		}

		var vault models.Vault
		err = db.First(&vault, vaultId).Error
		if err != nil {
			// We don't check if err is gorm.ErrRecordNotFound because we already checked it while checking
			// user's vault permissions.
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Retrieving vault from db failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, VaultRetrieveResponse{
			Id:        vault.ID,
			Name:      vault.Name,
			CreatedAt: vault.CreatedAt,
		})
	}
}

// HandleVaultDelete
//
//	@Summary	Delete vault by id
//	@Tags		vaults
//	@Id			deleteVault
//	@Success	204
//	@Failure	401
//	@Forbidden	403
//	@Failure	404
//	@Failure	500
//	@Router		/vaults/{id} [delete]
//	@Param		id	path	int	true	"Vault id"
func HandleVaultDelete(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		canDelete, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionDeleteVault)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canDelete {
			c.Status(http.StatusForbidden)
			return
		}

		err = db.Delete(&models.Vault{}, vaultId).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Deleting Vault failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Where("vault_id = ?", vaultId).Delete(&models.VaultPermission{}).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Deleting Vault Permissions failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Where("vault_id = ?", vaultId).Delete(&models.VaultKey{}).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Deleting vault keys failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Where("vault_id = ?", vaultId).Delete(&models.VaultItem{}).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Deleting vault items failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		auditLog := models.VaultAuditLog{
			VaultID:     uint(vaultId),
			VaultItemID: 0,
			UserID:      user.ID,
			ActionCode:  models.AuditLogActionVaultDelete,
			ActionData:  models.AuditLogDataVaultDelete(),
		}
		if err := db.Create(&auditLog).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Saving audit log failed.")
		}

		c.Status(http.StatusNoContent)
	}
}

// HandleVaultsMyPermissions
//
//	@Summary	List current user's permission on vault
//	@Tags		vaults
//	@Id			listMyVaultPermissions
//	@Produce	json
//	@Success	200	{object}	[]string
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	500
//	@Router		/vaults/{id}/my-permissions [get]
//	@Param		id	path	int	true	"Vault id"
func HandleVaultsMyPermissions(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		permissions := []string{}
		err = db.Model(&models.VaultPermission{}).Select("permission").
			Where("vault_id = ? AND user_id = ?", vaultId, user.ID).Scan(&permissions).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying vault permissions failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, permissions)
	}
}

// HandleVaultsLeave
//
//	@Summary	Leave from the vault
//	@Tags		vaults
//	@Id			leaveVault
//	@Produce	json
//	@Success	204
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/vaults/{id}/leave [post]
//	@Param		id	path	int	true	"Vault id"
func HandleVaultsLeave(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Where("vault_id = ? AND user_id = ?", vaultId, user.ID).Delete(&models.VaultPermission{}).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Removing user vault permissions failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Unscoped().Where("vault_id = ? AND key_owner_user_id = ?", vaultId, user.ID).Delete(&models.VaultKey{}).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Removing user vault key failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		auditLog := models.VaultAuditLog{
			VaultID:     uint(vaultId),
			VaultItemID: 0,
			UserID:      user.ID,
			ActionCode:  models.AuditLogActionVaultUserLeft,
			ActionData:  models.AuditLogDataVaultUserLeft(),
		}
		if err := db.Create(&auditLog).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Saving audit log failed.")
		}

		c.Status(http.StatusNoContent)
	}
}

// HandleVaultsMyKey
//
//	@Summary	Retrieve current user's vault key record for the vault
//	@Tags		vaults
//	@Id			retrieveMyVaultKey
//	@Produce	json
//	@Success	200	{object}	controllers.HandleVaultsMyKey.VaultKeyResponse
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/vaults/{id}/key [get]
//	@Param		id	path	int	true	"Vault id"
func HandleVaultsMyKey(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultKeyResponse struct {
		KeyOwnerUserID       int    `json:"key_owner_user_id" binding:"required"`
		InviterUserID        int    `json:"inviter_user_id" binding:"required"`
		EncryptionIV         string `json:"encryption_iv" binding:"required"`
		EncryptedVaultKey    string `json:"encrypted_vault_key" binding:"required"`
		InviterUserPublicKey string `json:"inviter_user_public_key" binding:"required"`
	}

	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		vaultKey := models.VaultKey{}
		err = db.Preload("InviterUser").First(&vaultKey, "vault_id = ? AND key_owner_user_id = ?", vaultId, user.ID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(http.StatusNotFound)
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying vault key of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, VaultKeyResponse{
			KeyOwnerUserID:       int(vaultKey.KeyOwnerUserID),
			InviterUserID:        int(vaultKey.InviterUserID),
			EncryptionIV:         vaultKey.EncryptionIV,
			EncryptedVaultKey:    vaultKey.EncryptedVaultKey,
			InviterUserPublicKey: vaultKey.InviterUser.PublicKey,
		})
	}
}

// HandleVaultsManageAddUser
//
//	@Summary	Add user to vault
//	@Tags		vault manage
//	@Id			addUserToVault
//	@Param		id		path	int														true	"Vault id"
//	@Param		request	body	controllers.HandleVaultsManageAddUser.AddUserRequest	true	"New user data"
//	@Produce	json
//	@Success	200
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults/{id}/manage/add-user [post]
func HandleVaultsManageAddUser(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type AddUserRequest struct {
		Email                string   `json:"email" binding:"required"`
		Permissions          []string `json:"permissions" binding:"required"`
		VaultKeyEncryptionIV string   `json:"vault_key_encryption_iv" binding:"required"`
		EncryptedVaultKey    string   `json:"encrypted_vault_key" binding:"required"`
	}

	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		canManageVault, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionManageVault)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canManageVault {
			c.Status(http.StatusForbidden)
			return
		}

		var requestData AddUserRequest
		if ok = bodybinder.Bind(&requestData, c); !ok {
			return
		}

		var newUser models.User
		err = db.First(&newUser, "email = ?", requestData.Email).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "User with given email not found."})
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying users by email failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		var isAlreadyAdded bool
		err = db.Model(&models.VaultPermission{}).Select("count(*) > 0").
			Where("vault_id = ? AND user_id = ?", vaultId, newUser.ID).Scan(&isAlreadyAdded).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying if user already added to vault failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if isAlreadyAdded {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "User is already added to vault."})
			return
		}

		// check and prepare vault permission record
		newUserVaultPermissions := []models.VaultPermission{}
		for _, p := range requestData.Permissions {
			if !slices.Contains([]string{
				models.VaultPermissionManageVault,
				models.VaultPermissionDeleteVault,
				models.VaultPermissionManageItems}, p) {
				c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: fmt.Sprintf("Given permission '%s' is invalid.", p)})
				return
			}
			newUserVaultPermissions = append(newUserVaultPermissions, models.VaultPermission{
				VaultID:    uint(vaultId),
				UserID:     newUser.ID,
				Permission: p,
			})
		}
		// append the read permission since it's mandatory
		newUserVaultPermissions = append(newUserVaultPermissions, models.VaultPermission{
			VaultID:    uint(vaultId),
			UserID:     newUser.ID,
			Permission: models.VaultPermissionRead,
		})

		err = db.Create(&newUserVaultPermissions).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating vault permissions for newly added user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		vaultKeyRecord := models.VaultKey{
			VaultID:           uint(vaultId),
			KeyOwnerUserID:    newUser.ID,
			InviterUserID:     user.ID,
			EncryptionIV:      requestData.VaultKeyEncryptionIV,
			EncryptedVaultKey: requestData.EncryptedVaultKey,
		}
		err = db.Create(&vaultKeyRecord).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating vault key record for newly added user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		auditLog := models.VaultAuditLog{
			VaultID:     uint(vaultId),
			VaultItemID: 0,
			UserID:      user.ID,
			ActionCode:  models.AuditLogActionVaultAddUser,
			ActionData: models.AuditLogDataVaultAddUser(newUser.Email, common.Map(newUserVaultPermissions,
				func(i models.VaultPermission) string { return i.Permission })),
		}
		if err := db.Create(&auditLog).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Saving audit log failed.")
		}

		c.Status(http.StatusOK)
	}
}

// HandleVaultsManageListUsers
//
//	@Summary	List users who have access to vault
//	@Tags		vault manage
//	@Id			listVaultUsers
//	@Param		id	path	int	true	"Vault id"
//	@Produce	json
//	@Success	200	{object}	[]controllers.HandleVaultsManageListUsers.UsersResponseItem
//	@Failure	401
//	@Failure	403
//	@Failure	500
//	@Router		/vaults/{id}/manage/users [get]
func HandleVaultsManageListUsers(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type UsersResponseItem struct {
		Id          int      `json:"id" binding:"required"`
		Email       string   `json:"email" binding:"required"`
		Permissions []string `json:"permissions" binding:"required"`
	}

	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		canManageVault, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionManageVault)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canManageVault {
			c.Status(http.StatusForbidden)
			return
		}

		var usersAndPermissions []struct {
			Id         int    `gorm:"column:id"`
			Email      string `gorm:"column:email"`
			Permission string `gorm:"column:permission"`
		}
		err = db.Select("users.id as id, users.email as email, vault_permissions.permission as permission").
			Model(&models.VaultPermission{}).
			Joins("LEFT OUTER JOIN users ON vault_permissions.user_id = users.id").
			Where("vault_permissions.vault_id = ?", vaultId).
			Order("permission ASC").
			Scan(&usersAndPermissions).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying users and permissions by vault failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		type UserKey struct {
			id    int
			email string
		}
		userAndPermissionsMap := make(map[UserKey][]string)
		for _, v := range usersAndPermissions {
			userKey := UserKey{v.Id, v.Email}
			_, ok := userAndPermissionsMap[userKey]
			if !ok {
				userAndPermissionsMap[userKey] = []string{}
			}
			userAndPermissionsMap[userKey] = append(userAndPermissionsMap[userKey], v.Permission)
		}

		result := []UsersResponseItem{}
		for k, v := range userAndPermissionsMap {
			result = append(result, UsersResponseItem{Id: k.id, Email: k.email, Permissions: v})
		}
		slices.SortFunc(result, func(i, j UsersResponseItem) int {
			return strings.Compare(i.Email, j.Email)
		})
		c.JSON(http.StatusOK, result)
	}
}

// HandleVaultsManageRemoveUser
//
//	@Summary	Remove user from vault
//	@Tags		vault manage
//	@Id			removeUserFromVault
//	@Param		id		path	int															true	"Vault id"
//	@Param		request	body	controllers.HandleVaultsManageRemoveUser.RemoveUserRequest	true	"ID of the user which will be removed"
//	@Success	204
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults/{id}/manage/users [delete]
func HandleVaultsManageRemoveUser(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type RemoveUserRequest struct {
		UserId int `json:"user_id" binding:"required"`
	}

	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		canManageVault, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionManageVault)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canManageVault {
			c.Status(http.StatusForbidden)
			return
		}

		var requestData RemoveUserRequest
		if ok = bodybinder.Bind(&requestData, c); !ok {
			return
		}

		var removedUser models.User
		err = db.First(&removedUser, requestData.UserId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(http.StatusNotFound)
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Where("vault_id = ? AND user_id = ?", vaultId, removedUser.ID).Delete(&models.VaultPermission{}).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Removing user from vault failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Unscoped().Where("vault_id = ? AND key_owner_user_id = ?", vaultId, removedUser.ID).Delete(&models.VaultKey{}).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Removing user from vault failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		auditLog := models.VaultAuditLog{
			VaultID:     uint(vaultId),
			VaultItemID: 0,
			UserID:      user.ID,
			ActionCode:  models.AuditLogActionVaultRemoveUser,
			ActionData:  models.AuditLogDataVaultRemoveUser(removedUser.Email),
		}
		if err := db.Create(&auditLog).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Saving audit log failed.")
		}

		c.Status(http.StatusNoContent)
	}
}

// HandleVaultsManageRename
//
//	@Summary	Rename vault
//	@Tags		vault manage
//	@Id			renameVault
//	@Param		id		path		int														true	"Vault id"
//	@Param		request	body		controllers.HandleVaultsManageRename.RenameVaultRequest	true	"New name of the vault"
//	@Success	200		{object}	controllers.HandleVaultsManageRename.RenameVaultResponse
//	@Failure	400		{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	404
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults/{id}/manage/rename [post]
func HandleVaultsManageRename(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type RenameVaultRequest struct {
		Name string `json:"name" binding:"required"`
	}

	type RenameVaultResponse struct {
		Name string `json:"name" binding:"required"`
	}

	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		canManageVault, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionManageVault)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canManageVault {
			c.Status(http.StatusForbidden)
			return
		}

		var requestData RenameVaultRequest
		if ok = bodybinder.Bind(&requestData, c); !ok {
			return
		}

		var vault models.Vault
		err = db.First(&vault, vaultId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(http.StatusNotFound)
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying vault failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		oldVaultName := vault.Name
		vault.Name = requestData.Name
		err = db.Save(&vault).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Saving vault's new name failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		auditLog := models.VaultAuditLog{
			VaultID:     vault.ID,
			VaultItemID: 0,
			UserID:      user.ID,
			ActionCode:  models.AuditLogActionVaultRename,
			ActionData:  models.AuditLogDataVaultRename(oldVaultName, vault.Name),
		}
		if err := db.Create(&auditLog).Error; err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Saving audit log failed.")
		}

		c.JSON(http.StatusOK, RenameVaultResponse{Name: requestData.Name})
	}
}

// HandleVaultAuditLogsList
//
//	@Summary	List audit logs of vault
//	@Tags		vaults
//	@Id			listVaultAuditLogs
//	@Produce	json
//	@Param		id			path		int	true	"Vault id"
//	@Param		page		query		int	false	"Page number"			default(1)	minimum(1)
//	@Param		page_size	query		int	false	"Item count per page"	default(10)
//	@Success	200			{object}	pagination.StandardPaginationResponse[controllers.HandleVaultAuditLogsList.AuditLogResponseItem]
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/vaults/{id}/logs [get]
func HandleVaultAuditLogsList(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type UserData struct {
		Id    uint   `json:"id" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	type VaultItemData struct {
		Id    uint   `json:"id" binding:"required"`
		Title string `json:"title" binding:"required"`
	}

	type AuditLogResponseItem struct {
		Id         uint                  `json:"id" binding:"required"`
		ActionCode models.AuditLogAction `json:"action_code" binding:"required"`
		ActionData map[string]any        `json:"action_data" binding:"required"`
		CreatedAt  time.Time             `json:"created_at" binding:"required"`
		User       UserData              `json:"user" binding:"required"`
		VaultItem  *VaultItemData        `json:"vault_item"`
	}

	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		canReadVault, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionRead)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canReadVault {
			c.Status(http.StatusForbidden)
			return
		}

		var count int64
		err = db.Model(&models.VaultAuditLog{}).Select("count(*)").
			Where("vault_id = ?", vaultId).Count(&count).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying audit log count failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		auditLogs := []models.VaultAuditLog{}
		err = db.Unscoped().Scopes(pagination.Paginate(c)).Joins("VaultItem").Joins("User").
			Order("vault_audit_logs.created_at DESC").Find(&auditLogs, "vault_audit_logs.vault_id = ?", vaultId).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying audit logs failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		results := make([]AuditLogResponseItem, len(auditLogs))
		for i, v := range auditLogs {
			log := AuditLogResponseItem{
				Id:         v.ID,
				ActionCode: v.ActionCode,
				ActionData: v.ActionData,
				CreatedAt:  v.CreatedAt,
				User: UserData{
					Id:    v.User.ID,
					Email: v.User.Email,
				},
			}
			if v.VaultItem.ID != 0 {
				log.VaultItem = &VaultItemData{
					Id:    v.VaultItem.ID,
					Title: v.VaultItem.Title,
				}
			}
			results[i] = log
		}

		c.JSON(http.StatusOK, pagination.StandardPaginationResponse[AuditLogResponseItem]{
			Results: results,
			Count:   int(count),
		})
	}
}
