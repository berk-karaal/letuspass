package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	postgresDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = postgresDb.AutoMigrate(GetModels()...)
	if err != nil {
		return nil, err
	}

	return postgresDb, nil
}

func GenerateDSN(host, user, password, dbname, port, sslmode, timeZone string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timeZone)
}
