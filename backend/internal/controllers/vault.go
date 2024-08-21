package controllers

import (
	"net/http"

	"github.com/berk-karaal/letuspass/backend/internal/common/bodybinder"
	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/models"
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
