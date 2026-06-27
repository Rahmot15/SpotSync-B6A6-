package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg Config) (*gorm.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database Connected")
	return db, nil
}

func AutoMigrate(db *gorm.DB, models ...any) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	return nil
}
