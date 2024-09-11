package main

import (
	"log"
	"os"

	"github.com/berk-karaal/letuspass/backend/internal/config"
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

	log.Printf("Starting rest api server.")
	router.SetupRouter(apiConfig).Run()
}
