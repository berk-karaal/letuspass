package controllers

import (
	"net/http"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// HandleUsersMe
//
//	@Summary	Get currently logged-in user
//	@Tags		users
//	@Produce	json
//	@Success	200	{object}	controllers.HandleUsersMe.MeResponse
//	@Failure	401
//	@Failure	500
//	@Router		/users/me [get]
func HandleUsersMe(logger *logging.Logger) func(c *gin.Context) {
	type MeResponse struct {
		Email string `json:"email"`
		Name  string `json:"name"`
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
