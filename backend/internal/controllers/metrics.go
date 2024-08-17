package controllers

import (
	"net/http"

	"github.com/berk-karaal/letuspass/backend/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// HandleMetricsStatus
//
//	@Summary	Get status of the server
//	@Tags		metrics
//	@Produce	json
//	@Success	200	{object}	controllers.HandleMetricsStatus.MetricsStatusResponse
//	@Router		/metrics/status [get]
func HandleMetricsStatus(logger *logging.Logger) func(c *gin.Context) {
	type MetricsStatusResponse struct {
		Status string `json:"status"`
	}

	return func(c *gin.Context) {
		// this log is made for demo purposes
		logger.RequestEvent(zerolog.InfoLevel, c).Msg("Status request received.")

		c.JSON(http.StatusOK, MetricsStatusResponse{Status: "OK"})
	}
}
