package routes

import (
	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(engine *gin.Engine, logger *logging.Logger, postgres *gorm.DB) {
	v1Group := engine.Group("/api/v1")
	{
		metricGroup := v1Group.Group("/metrics")
		{
			metricGroup.GET("/status", controllers.HandleMetricsStatus(logger))
		}

		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/login", controllers.HandleAuthLogin(logger, postgres))
			authGroup.POST("/register", controllers.HandleAuthRegister(logger, postgres))
		}
	}
}
