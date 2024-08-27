package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/berk-karaal/letuspass/backend/internal/common/bodybinder"
	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/config"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/models"
	"github.com/berk-karaal/letuspass/backend/internal/schemas"
	authservice "github.com/berk-karaal/letuspass/backend/internal/services/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// HandleAuthLogin
//
//	@Summary	Login user
//	@Tags		auth
//	@Id			authLogin
//	@Param		request	body	controllers.HandleAuthLogin.LoginRequest	true	"Login credentials"
//	@Produce	json
//	@Success	200	{object}	controllers.HandleAuthLogin.LoginResponse
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/auth/login [post]
func HandleAuthLogin(apiConfig *config.RestapiConfig, logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	type LoginResponse struct {
		Email string `json:"email" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	return func(c *gin.Context) {
		isAlreadyAuthenticated := true
		_, _, err := middlewares.GetCurrentUser(c, apiConfig, db)
		if err != nil {
			if errors.Is(err, middlewares.UserNotAuthenticatedErr{}) {
				isAlreadyAuthenticated = false
			} else {
				logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking if user already authenticated failed")
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		if isAlreadyAuthenticated {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "User already logged-in."})
			return
		}

		var requestData LoginRequest
		if !bodybinder.Bind(&requestData, c) {
			return
		}

		var user models.User
		err = db.First(&user, "email = ?", requestData.Email).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Wrong credentials."})
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Getting user by email failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		ok, err := authservice.ComparePassword(user.Password, requestData.Password)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Comparing password failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Wrong credentials."})
			return
		}

		session := models.UserSession{
			Token:     authservice.GenerateSessionToken(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Second * time.Duration(apiConfig.SessionTokenExpireSeconds)),
			UserAgent: c.Request.UserAgent(),
		}

		err = db.Create(&session).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating user session failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.SetCookie(apiConfig.SessionTokenCookieName, session.Token, apiConfig.SessionTokenExpireSeconds, "/", "localhost", true, true)

		c.JSON(http.StatusOK, LoginResponse{Email: user.Email, Name: user.Name})
	}
}

// HandleAuthRegister
//
//	@Summary	Register user
//	@Tags		auth
//	@Id			authRegister
//	@Param		request	body	controllers.HandleAuthRegister.RegisterRequest	true	"User Registration Data"
//	@Produce	json
//	@Success	201
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Failure	500
//	@Router		/auth/register [post]
func HandleAuthRegister(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type RegisterRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	return func(c *gin.Context) {
		var requestData RegisterRequest
		if !bodybinder.Bind(&requestData, c) {
			return
		}

		var exists bool
		err := db.Model(&models.User{}).Select("count(*) > 0").Where("email = ?", requestData.Email).Scan(&exists).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking if an user with same email exists failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "User with given email already exists."})
			return
		}

		hashedPassword, err := authservice.HashPassword(requestData.Password)
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Hashing password failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		newUser := models.User{
			Email:    requestData.Email,
			Password: hashedPassword,
			Name:     requestData.Name,
			IsActive: true,
		}
		err = db.Create(&newUser).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating new user failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusCreated)
	}
}

// HandleAuthLogout
//
//	@Summary	Logout user
//	@Tags		auth
//	@Id			authLogout
//	@Produce	json
//	@Success	204
//	@Failure	401
//	@Failure	500
//	@Router		/auth/logout [post]
func HandleAuthLogout(apiConfig *config.RestapiConfig, logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie(apiConfig.SessionTokenCookieName)

		tx := db.Delete(&models.UserSession{}, "token = ?", sessionToken)
		if tx.Error != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Deleting UserSession failed.")
			c.Status(http.StatusInternalServerError)
			return
		}
		if tx.RowsAffected != 1 {
			logger.RequestEvent(zerolog.ErrorLevel, c).Str("used_token", sessionToken).
				Msg(fmt.Sprintf("Deleted %d UserSession rows instead of just 1.", tx.RowsAffected))
		}

		// Remove cookie
		c.SetCookie(apiConfig.SessionTokenCookieName, "", 0, "/", "localhost", true, true)

		c.Status(http.StatusNoContent)
	}
}
