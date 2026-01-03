package database

import (
	"learning-app-backend/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) {
	var err error
	var dialector gorm.Dialector

	// Choose database driver based on configuration
	if cfg.DatabaseType == "postgres" {
		log.Println("Connecting to PostgreSQL database...")
		dialector = postgres.Open(cfg.DatabaseURL)
	} else {
		log.Println("Connecting to SQLite database...")
		dialector = sqlite.Open(cfg.DatabaseURL)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Printf("Database connection established (%s)", cfg.DatabaseType)
}
