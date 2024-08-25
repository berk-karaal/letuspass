package controllers

import (
	"errors"
	"net/http"
	"strconv"

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
//	@Router		/vaults/{id}/items [post]
//	@Param		id	path	int	true	"Vault id"
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

// HandleVaultItemsList
//
//	@Summary	List items of a vault
//	@Tags		vault items
//	@Param		page		query	int		false	"Page number"			default(1)	minimum(1)
//	@Param		page_size	query	int		false	"Item count per page"	default(10)
//	@Param		ordering	query	string	false	"Ordering"				Enums(title, -title, created_at, -created_at)
//	@Produce	json
//	@Success	200	{object}	pagination.StandardPaginationResponse[controllers.HandleVaultItemsList.VaultItemResponseItem]
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	500
//	@Router		/vaults/{id}/items [get]
//	@Param		id	path	int	true	"Vault id"
func HandleVaultItemsList(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultItemResponseItem struct {
		Id    uint   `json:"id"`
		Title string `json:"title"`
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

		ordering, err := orderbyparam.GenerateOrdering(c, map[string]string{
			"title":      "vault_items.title",
			"created_at": "vault_items.created_at",
		}, "title")
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Generating query ordering from params failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		var count int64
		var results []VaultItemResponseItem
		err = db.Scopes(pagination.Paginate(c)).Select("id, title").Table("vault_items").
			Where("deleted_at IS NULL AND vault_id = ?", vaultId).Order(ordering).
			Count(&count).Scan(&results).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying vault items failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, pagination.StandardPaginationResponse[VaultItemResponseItem]{
			Count:   int(count),
			Results: results,
		})
	}
}

// HandleVaultItemsRetrieve
//
//	@Summary	Retrieve a new vault item
//	@Tags		vault items
//	@Success	200	{object}	controllers.HandleVaultItemsRetrieve.VaultItemRetrieveResponse
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	404	{object}	schemas.NotFoundResponse
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults/{id}/items/{itemId} [get]
//	@Param		id		path	int	true	"Vault id"
//	@Param		itemId	path	int	true	"Vault Item id"
func HandleVaultItemsRetrieve(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultItemRetrieveResponse struct {
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

		vaultItemId, err := strconv.Atoi(c.Param("itemId"))
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

		var vaultItem models.VaultItem
		err = db.First(&vaultItem, "id = ? AND vault_id = ?", vaultItemId, vaultId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, schemas.NotFoundResponse{Error: "Vault item doesn't exist."})
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Getting vault item from database failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, VaultItemRetrieveResponse{
			Id:                vaultItem.ID,
			Title:             vaultItem.Title,
			EncryptedUsername: vaultItem.EncryptedUsername,
			EncryptedPassword: vaultItem.EncryptedPassword,
			EncryptedNote:     vaultItem.EncryptedNote,
		})
	}
}

// HandleVaultItemsUpdate
//
//	@Summary	Update a new vault item
//	@Tags		vault items
//	@Param		request	body	controllers.HandleVaultItemsUpdate.VaultItemUpdateRequest	true	"New vault item data"
//	@Produce	json
//	@Success	200	{object}	controllers.HandleVaultItemsUpdate.VaultItemUpdateResponse
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	404	{object}	schemas.NotFoundResponse
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults/{id}/items/{itemId} [put]
//	@Param		id		path	int	true	"Vault id"
//	@Param		itemId	path	int	true	"Vault Item id"
func HandleVaultItemsUpdate(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type VaultItemUpdateRequest struct {
		Title             string `json:"title" binding:"required"`
		EncryptedUsername string `json:"encrypted_username"`
		EncryptedPassword string `json:"encrypted_password"`
		EncryptedNote     string `json:"encrypted_note"`
	}

	type VaultItemUpdateResponse struct {
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

		vaultItemId, err := strconv.Atoi(c.Param("itemId"))
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

		var vaultItem models.VaultItem
		err = db.First(&vaultItem, vaultItemId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, schemas.NotFoundResponse{Error: "Vault item doesn't exist."})
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Getting vault item from database failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if int(vaultItem.VaultID) != vaultId {
			c.JSON(http.StatusNotFound, schemas.NotFoundResponse{Error: "Vault item doesn't exist."})
			return
		}

		var requestData VaultItemUpdateRequest
		if ok = bodybinder.Bind(&requestData, c); !ok {
			return
		}

		vaultItem.Title = requestData.Title
		vaultItem.EncryptedUsername = requestData.EncryptedUsername
		vaultItem.EncryptedPassword = requestData.EncryptedPassword
		vaultItem.EncryptedNote = requestData.EncryptedNote
		err = db.Save(&vaultItem).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Updating vault item failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		// TODO: audit log

		c.JSON(http.StatusOK, VaultItemUpdateResponse{
			Id:                vaultItem.ID,
			Title:             vaultItem.Title,
			EncryptedUsername: vaultItem.EncryptedUsername,
			EncryptedPassword: vaultItem.EncryptedPassword,
			EncryptedNote:     vaultItem.EncryptedNote,
		})
	}
}

// HandleVaultItemsDelete
//
//	@Summary	Delete a vault item
//	@Tags		vault items
//	@Success	204
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	401
//	@Failure	403
//	@Failure	404	{object}	schemas.NotFoundResponse
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/vaults/{id}/items/{itemId} [delete]
//	@Param		id		path	int	true	"Vault id"
//	@Param		itemId	path	int	true	"Vault Item id"
func HandleVaultItemsDelete(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		vaultId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Id must be an integer."})
			return
		}

		vaultItemId, err := strconv.Atoi(c.Param("itemId"))
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

		var vaultItem models.VaultItem
		err = db.First(&vaultItem, "id = ? AND vault_id = ?", vaultItemId, vaultId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, schemas.NotFoundResponse{Error: "Vault item doesn't exist."})
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Getting vault item from database failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		err = db.Unscoped().Delete(&vaultItem).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Deleting vault item failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		// TODO: audit log

		c.Status(http.StatusNoContent)
	}
}
