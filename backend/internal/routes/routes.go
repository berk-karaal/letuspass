package routes

import (
	"github.com/berk-karaal/letuspass/backend/internal/controllers"
	"github.com/berk-karaal/letuspass/backend/internal/logging"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, logger *logging.Logger) {
	v1Group := engine.Group("/api/v1")
	{
		metricGroup := v1Group.Group("/metrics")
		{
			metricGroup.GET("/status", controllers.HandleMetricsStatus(logger))
		}
	}
}
