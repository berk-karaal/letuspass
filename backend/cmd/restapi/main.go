package main

import (
	"reflect"
	"strings"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/routes"
	_ "github.com/berk-karaal/letuspass/backend/swagger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// main
//
//	@title			LetusPass REST API
//	@version		0.0.1
//	@description	Project description at https://github.com/berk-karaal/letuspass
//	@host			localhost:8080
//	@BasePath		/api/v1
func main() {
	registerJsonTagNames()

	logger := logging.NewLogger()

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(requestid.New())
	engine.Use(middlewares.LogHandler(logger))
	routes.SetupRoutes(engine, logger)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.Run()
}

// registerJsonTagNames registers json tag names to gin validator. This registration is
// necessary to get json tag names of fields when validation fails. We need json tag
// names of fields when returning a response if validation fails.
//
// Ref:
//   - https://blog.depa.do/post/gin-validation-errors-handling#toc_6
//   - https://github.com/go-playground/validator/issues/258#issuecomment-257281334
func registerJsonTagNames() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
