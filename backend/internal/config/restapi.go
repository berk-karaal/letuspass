package config

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// use https://github.com/joho/godotenv

type RestapiConfig struct {
	GinMode string
	LogFile string

	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
	DbSSLMode  string
	DbTimeZone string

	SessionTokenCookieName    string
	SessionTokenExpireSeconds int
}

// NewRestapiConfigFromEnv creates a RestapiConfig from environment variables. It panics if converting types
// fails.
func NewRestapiConfigFromEnv() RestapiConfig {
	sessionTokenExpireSeconds, err := strconv.Atoi(os.Getenv("SESSION_TOKEN_EXPIRE_SECONDS"))
	if err != nil {
		log.Fatal("SESSION_TOKEN_EXPIRE_SECONDS env must be a valid integer")
	}

	var ginMode string
	switch os.Getenv("GIN_MODE") {
	case "debug":
		ginMode = gin.DebugMode
	case "test":
		ginMode = gin.TestMode
	case "release":
		ginMode = gin.ReleaseMode
	default:
		log.Fatal("GIN_MODE must be one of these: debug, test, release.")
	}

	return RestapiConfig{
		GinMode: ginMode,
		LogFile: os.Getenv("LOG_FILE"),

		DbHost:     os.Getenv("DB_HOST"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbPort:     os.Getenv("DB_PORT"),
		DbSSLMode:  os.Getenv("DB_SSL_MODE"),
		DbTimeZone: os.Getenv("DB_TIME_ZONE"),

		SessionTokenCookieName:    os.Getenv("SESSION_TOKEN_COOKIE_NAME"),
		SessionTokenExpireSeconds: sessionTokenExpireSeconds,
	}
}
