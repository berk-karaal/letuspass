package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleMetricsStatus() func(c *gin.Context) {
	type MetricsStatusResponse struct {
		Status string `json:"status"`
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, MetricsStatusResponse{Status: "OK"})
	}
}
