package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/berk-karaal/letuspass/backend/internal/common/bodybinder"
	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
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
//	@Param		request	body	controllers.HandleAuthLogin.LoginRequest	true	"Login credentials"
//	@Produce	json
//	@Success	200
//	@Failure	400	{object}	schemas.BadRequestResponse
//	@Failure	422	{object}	bodybinder.validationErrorResponse
//	@Router		/auth/login [post]
func HandleAuthLogin(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	type LoginResponse struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	return func(c *gin.Context) {
		var requestData LoginRequest
		if !bodybinder.Bind(&requestData, c) {
			return
		}

		var user models.User
		err := db.First(&user, "email = ?", requestData.Email).Error
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
			ExpiresAt: time.Time{},
			UserAgent: c.Request.UserAgent(),
		}

		err = db.Create(&session).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Creating user session failed.")
			c.Status(http.StatusInternalServerError)
			return
		}

		// TODO: get from config
		c.SetCookie("session_token", session.Token, int((time.Hour * 24).Seconds()), "/", "localhost", true, true)

		c.JSON(http.StatusOK, LoginResponse{Email: user.Email, Name: user.Name})
	}
}

// HandleAuthRegister
//
//	@Summary	Register user
//	@Tags		auth
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
