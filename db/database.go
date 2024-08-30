package db

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/vector-ops/ticketeer/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(config *config.EnvConfig, dbMigrator func(db *gorm.DB) error) *gorm.DB {
	uri := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s port=5432",
		config.DBHost, config.DBUser, config.DBName, config.DBPassword, config.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database %e", err)
	}

	log.Info("Connected to database")

	if err := dbMigrator(db); err != nil {
		log.Fatalf("Failed to migrate tables %e", err)
	}
	return db
}
