package controllers

import (
	"errors"
	"net/http"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// HandleUsersMe
//
//	@Summary	Get currently logged-in user
//	@Tags		users
//	@Id			getCurrentUser
//	@Produce	json
//	@Success	200	{object}	controllers.HandleUsersMe.MeResponse
//	@Failure	401
//	@Failure	500
//	@Router		/users/me [get]
func HandleUsersMe(logger *logging.Logger) func(c *gin.Context) {
	type MeResponse struct {
		Email string `json:"email" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	return func(c *gin.Context) {
		user, ok := middlewares.ExtractUserFromGinContext(c)
		if !ok {
			logger.RequestEvent(zerolog.ErrorLevel, c).Msg("Extracting user from Gin context failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, MeResponse{
			Email: user.Email,
			Name:  user.Name,
		})
	}
}

// HandleGetUserByEmail
//
//	@Summary	Get user by email
//	@Tags		users
//	@Id			getUserByEmail
//	@Produce	json
//	@Success	200	{object}	controllers.HandleGetUserByEmail.UserResponse
//	@Failure	401
//	@Failure	404
//	@Failure	500
//	@Router		/users/by-email [get]
//	@Param		email	query	string	true	"Email of the user"
func HandleGetUserByEmail(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type UserResponse struct {
		Email     string `json:"email" binding:"required"`
		Name      string `json:"name" binding:"required"`
		PublicKey string `json:"public_key" binding:"required"`
	}

	return func(c *gin.Context) {
		email := c.Query("email")

		var user models.User
		err := db.First(&user, "email = ?", email).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(http.StatusNotFound)
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying user by email failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			Email:     user.Email,
			Name:      user.Name,
			PublicKey: user.PublicKey,
		})
	}
}
