package main

import (
	"github.com/berk-karaal/letuspass/backend/internal/logging"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/routes"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.NewLogger()

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(requestid.New())
	engine.Use(middlewares.LogHandler(logger))
	routes.SetupRoutes(engine)

	engine.Run()
}
