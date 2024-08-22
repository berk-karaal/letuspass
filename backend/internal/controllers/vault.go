package controllers

import (
	"net/http"
	"strconv"
	"time"

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
//	@Param		request	body	controllers.HandleVaultsCreate.VaultCreateRequest	true	"New vault data"
//	@Produce	json
//	@Success	201	{object}	controllers.HandleVaultsCreate.VaultCreateResponse
//	@Failure	401
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults [post]
func HandleVaultsCreate(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultCreateRequest struct {
		Name string `json:"name" binding:"required"`
	}

	type VaultCreateResponse struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
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

		// TODO: audit log

		c.JSON(http.StatusCreated, VaultCreateResponse{Id: vault.ID, Name: vault.Name})
	}
}

// HandleVaultsList
//
//	@Summary	List vaults that user has read access to
//	@Tags		vaults
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
		Id        uint      `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
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

		var results []VaultResponseItem
		var count int64
		err = db.Scopes(pagination.Paginate(c)).Select("vaults.id, vaults.name, vaults.created_at, vaults.updated_at").
			Table("vault_permissions").
			Joins("LEFT OUTER JOIN vaults ON vault_permissions.vault_id = vaults.id").
			Where("vault_permissions.deleted_at IS NULL AND vault_permissions.user_id = ? AND vault_permissions.permission = ?",
				user.ID, models.VaultPermissionRead).
			Order(ordering).
			Count(&count).
			Scan(&results).Error
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
		Id        uint      `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
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

		// TODO: delete vault related data
		// Delete vault items
		// Delete vault keys
		// ...

		c.Status(http.StatusNoContent)
	}
}
