package main

import (
	"github.com/berk-karaal/letuspass/backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	routes.SetupRoutes(engine)
	engine.Run()
}
