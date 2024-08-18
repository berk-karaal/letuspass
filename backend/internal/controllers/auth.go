package controllers

import (
	"github.com/berk-karaal/letuspass/backend/internal/common/bodybinder"
	"github.com/gin-gonic/gin"
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
func HandleAuthLogin() func(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var requestData LoginRequest

		if !bodybinder.Bind(&requestData, c) {
			return
		}

		// TODO
	}
}
