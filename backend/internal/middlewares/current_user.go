package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const UserContextKey = "user"

// CurrentUserHandler middleware checks request for the authenticated user and puts user data to Gin context.
// If request sent without authentication, this middleware aborts with HTTP 401 Unauthorized with no response
// body. This middleware should only be used on routes that required authentication.
func CurrentUserHandler(logger *logging.Logger, db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token") // TODO: get from config
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Checking session cookie failed.")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var userSession models.UserSession
		err = db.First(&userSession, "token = ?", sessionToken).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying session token failed.")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = db.Model(&userSession).Update("expires_at", time.Now().Add(time.Minute*60*24)).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Extending user session expire time failed.")
		}

		var user models.User
		err = db.First(&user, "id = ?", userSession.UserID).Error
		if err != nil {
			logger.RequestEvent(zerolog.ErrorLevel, c).Err(err).Msg("Querying user failed.")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set(UserContextKey, user)
	}
}

// ExtractUserFromGinContext returns User model inserted to Gin context by CurrentUserHandler middleware.
// This helper function is used to reduce boilerplate code.
func ExtractUserFromGinContext(c *gin.Context) (models.User, bool) {
	val, ok := c.Get(UserContextKey)
	if !ok {
		return models.User{}, false
	}
	user, ok := val.(models.User)
	if !ok {
		return models.User{}, false
	}
	return user, true
}

type UserNotAuthenticatedErr struct{}

func (e UserNotAuthenticatedErr) Error() string { return "user is not authenticated" }

// GetCurrentUser returns the User and UserSession if the user is logged-in, if not returns
// UserNotAuthenticatedErr.
func GetCurrentUser(c *gin.Context, db *gorm.DB) (models.User, models.UserSession, error) {
	sessionToken, err := c.Cookie("session_token") // TODO: get from config
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return models.User{}, models.UserSession{}, UserNotAuthenticatedErr{}
		}
		return models.User{}, models.UserSession{}, fmt.Errorf("getting session_token cookie failed: %w", err)
	}

	var userSession models.UserSession
	err = db.First(&userSession, "token = ?", sessionToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, models.UserSession{}, UserNotAuthenticatedErr{}
		}
		return models.User{}, models.UserSession{}, fmt.Errorf("querying UserSession by token failed: %w", err)
	}

	var user models.User
	err = db.First(&user, "id = ?", userSession.UserID).Error
	if err != nil {
		return models.User{}, models.UserSession{}, fmt.Errorf("querying User by id failed: %w", err)
	}

	return user, userSession, nil
}
