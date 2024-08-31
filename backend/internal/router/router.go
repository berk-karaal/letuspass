package router

import (
	"fmt"
	golog "log"
	"reflect"
	"strings"
	"time"

	"github.com/berk-karaal/letuspass/backend/internal/common/logging"
	"github.com/berk-karaal/letuspass/backend/internal/config"
	"github.com/berk-karaal/letuspass/backend/internal/databases/postgres"
	"github.com/berk-karaal/letuspass/backend/internal/middlewares"
	"github.com/berk-karaal/letuspass/backend/internal/models"
	_ "github.com/berk-karaal/letuspass/backend/swagger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(apiConfig config.RestapiConfig) *gin.Engine {
	registerJsonTagNames()

	logger := logging.NewLogger(apiConfig.LogFile)

	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		apiConfig.DbHost, apiConfig.DbUser, apiConfig.DbPassword, apiConfig.DbName,
		apiConfig.DbPort, apiConfig.DbSSLMode, apiConfig.DbTimeZone)
	postgresDb, err := postgres.NewDB(postgresDsn)
	if err != nil {
		golog.Fatal(err)
	}
	err = postgresDb.AutoMigrate(&models.User{}, &models.UserSession{}, &models.Vault{}, &models.VaultPermission{},
		&models.VaultItem{})
	if err != nil {
		golog.Fatal(err)
	}

	gin.SetMode(apiConfig.GinMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(middlewares.LogHandler(logger))
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.1.107:5173"},
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
