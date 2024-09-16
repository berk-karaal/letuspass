package main

import (
	"log"
	"os"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/config"
	"github.com/berk-karaal/letuspass/backend/internal/databases/postgres"
	"github.com/berk-karaal/letuspass/backend/internal/router"
	"github.com/joho/godotenv"
)

// main
//
//	@title			LetusPass REST API
//	@version		0.0.1
//	@description	Project description at https://github.com/berk-karaal/letuspass
//	@host			localhost:8080
//	@BasePath		/api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
	apiConfig := config.NewRestapiConfigFromEnv()

	logger := logging.NewLogger(apiConfig.LogFile)

	postgresDb, err := postgres.NewDB(postgres.GenerateDSN(
		apiConfig.DbHost, apiConfig.DbUser, apiConfig.DbPassword, apiConfig.DbName,
		apiConfig.DbPort, apiConfig.DbSSLMode, apiConfig.DbTimeZone))
	if err != nil {
		log.Fatal("Failed to connect to postgres database: ", err)
	}

	log.Printf("Starting rest api server.")
	router.SetupRouter(apiConfig, logger, postgresDb).Run()
}
