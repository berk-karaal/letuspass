package tests

import (
	"log"

	"github.com/berk-karaal/letuspass/backend/internal/databases/postgres"
	"gorm.io/gorm"
)

func CleanDatabase(db *gorm.DB) {
	err := db.Migrator().DropTable(postgres.GetModels()...)
	if err != nil {
		log.Fatalf("failed to drop tables: %v", err)
	}
}
