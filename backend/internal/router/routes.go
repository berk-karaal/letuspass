package router

import (
	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/config"
	"github.com/berk-karaal/letuspass/backend/internal/controllers"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(engine *gin.Engine, apiConfig *config.RestapiConfig, logger *logging.Logger, postgres *gorm.DB) {
	v1Group := engine.Group("/api/v1")
	{
		metricGroup := v1Group.Group("/metrics")
		{
			metricGroup.GET("/status", controllers.HandleMetricsStatus(logger))
		}

		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/login", controllers.HandleAuthLogin(apiConfig, logger, postgres))
			authGroup.POST("/register", controllers.HandleAuthRegister(logger, postgres))
			authGroup.POST("/logout", middlewares.CurrentUserHandler(apiConfig, logger, postgres), controllers.HandleAuthLogout(apiConfig, logger, postgres))
		}

		userGroup := v1Group.Group("/users", middlewares.CurrentUserHandler(apiConfig, logger, postgres))
		{
			userGroup.GET("/me", controllers.HandleUsersMe(logger))
		}

		vaultGroup := v1Group.Group("/vaults", middlewares.CurrentUserHandler(apiConfig, logger, postgres))
		{
			vaultGroup.POST("", controllers.HandleVaultsCreate(logger, postgres))
			vaultGroup.GET("", controllers.HandleVaultsList(logger, postgres))
			vaultGroup.GET("/:id", controllers.HandleVaultsRetrieve(logger, postgres))
			vaultGroup.DELETE("/:id", controllers.HandleVaultDelete(logger, postgres))

			vaultManage := vaultGroup.Group("/:id/manage")
			{
				vaultManage.GET("/users", controllers.HandleVaultsManageListUsers(logger, postgres))
				vaultManage.DELETE("/users", controllers.HandleVaultsManageRemoveUser(logger, postgres))
				vaultManage.POST("/add-user", controllers.HandleVaultsManageAddUser(logger, postgres))
			}
		}
	}
}
