package router

import (
	"reflect"
	"strings"
	"time"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/config"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	_ "github.com/berk-karaal/letuspass/backend/swagger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(apiConfig config.RestapiConfig, logger *logging.Logger, postgresDb *gorm.DB) *gin.Engine {
	registerJsonTagNames()

	gin.SetMode(apiConfig.GinMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(middlewares.LogHandler(logger))
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowOrigins:     apiConfig.CORSAllowOrigins,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	SetupRoutes(router, &apiConfig, logger, postgresDb)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
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
