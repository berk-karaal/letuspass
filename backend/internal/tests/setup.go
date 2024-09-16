package tests

import (
	"log"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/config"
	"github.com/berk-karaal/letuspass/backend/internal/databases/postgres"
	"github.com/berk-karaal/letuspass/backend/internal/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTestRouter() (r *gin.Engine, apiConfig config.RestapiConfig, postgresDb *gorm.DB) {
	apiConfig = config.NewRestapiConfigFromEnv()
	logger := logging.NewLogger(apiConfig.LogFile)

	postgresDb, err := postgres.NewDB(postgres.GenerateDSN(
		apiConfig.DbHost, apiConfig.DbUser, apiConfig.DbPassword, apiConfig.DbName,
		apiConfig.DbPort, apiConfig.DbSSLMode, apiConfig.DbTimeZone))
	if err != nil {
		log.Fatal("Failed to connect to postgres database: ", err)
	}

	return router.SetupRouter(apiConfig, logger, postgresDb), apiConfig, postgresDb
}
