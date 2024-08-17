package main

import (
	"github.com/berk-karaal/letuspass/backend/internal/logging"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/routes"
	_ "github.com/berk-karaal/letuspass/backend/swagger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			LetusPass REST API
//	@version		0.0.1
//	@description	Project description at https://github.com/berk-karaal/letuspass
//	@host			localhost:8080
//	@BasePath		/api/v1
func main() {
	logger := logging.NewLogger()

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(requestid.New())
	engine.Use(middlewares.LogHandler(logger))
	routes.SetupRoutes(engine, logger)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.Run()
}
