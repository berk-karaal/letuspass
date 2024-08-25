package controllers

import (
	"net/http"
	"strconv"

	"github.com/berk-karaal/letuspass/backend/internal/common/bodybinder"
	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/models"
	"github.com/berk-karaal/letuspass/backend/internal/schemas"
	vaultservice "github.com/berk-karaal/letuspass/backend/internal/services/vault"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// HandleVaultItemsCreate
//
//	@Summary	Create a new vault item
//	@Tags		vault items
//	@Param		request	body	controllers.HandleVaultItemsCreate.VaultItemCreateRequest	true	"New vault item data"
//	@Produce	json
//	@Success	201	{object}	controllers.HandleVaultItemsCreate.VaultItemCreateResponse
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults/:id/items [post]
func HandleVaultItemsCreate(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultItemCreateRequest struct {
		Title             string `json:"title" binding:"required"`
		EncryptedUsername string `json:"encrypted_username"`
		EncryptedPassword string `json:"encrypted_password"`
		EncryptedNote     string `json:"encrypted_note"`
	}

	type VaultItemCreateResponse struct {
		Id                uint   `json:"id"`
		Title             string `json:"title"`
		EncryptedUsername string `json:"encrypted_username"`
		EncryptedPassword string `json:"encrypted_password"`
		EncryptedNote     string `json:"encrypted_note"`
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

		canManageItems, err := vaultservice.CheckUserHasVaultPermission(db, int(user.ID), vaultId, models.VaultPermissionManageItems)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking vault permissions of user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !canManageItems {
			c.Status(http.StatusForbidden)
			return
		}

		var vault models.Vault
		err = db.First(&vault, vaultId).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Retrieving vault from db failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		var requestData VaultItemCreateRequest
		if ok = bodybinder.Bind(&requestData, c); !ok {
			return
		}

		vaultItem := models.VaultItem{
			VaultID:           uint(vaultId),
			Title:             requestData.Title,
			EncryptedUsername: requestData.EncryptedUsername,
			EncryptedPassword: requestData.EncryptedPassword,
			EncryptedNote:     requestData.EncryptedNote,
		}
		err = db.Create(&vaultItem).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating vault item failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		// TODO: audit log

		c.JSON(http.StatusCreated, VaultItemCreateResponse{
			Id:                vaultItem.ID,
			Title:             vaultItem.Title,
			EncryptedUsername: vaultItem.EncryptedUsername,
			EncryptedPassword: vaultItem.EncryptedPassword,
			EncryptedNote:     vaultItem.EncryptedNote,
		})
	}
}
