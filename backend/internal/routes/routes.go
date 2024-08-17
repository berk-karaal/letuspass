package routes

import (
	"github.com/berk-karaal/letuspass/backend/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine) {
	v1Group := engine.Group("/api/v1")
	{
		metricGroup := v1Group.Group("/metrics")
		{
			metricGroup.GET("/status", controllers.HandleMetricsStatus())
		}
	}
}
