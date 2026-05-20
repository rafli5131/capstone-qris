package config

import (
	"fmt"
	"time"

	"capstone-qris/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewDatabase(cfg *Config, log *logrus.Logger) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Fatal("failed to get underlying sql.DB")
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// Auto-migrate GORM models (no-op when using golang-migrate; safe to keep for dev)
	_ = db.AutoMigrate(&entity.ApiClient{}, &entity.Merchant{}, &entity.Account{}, &entity.Transaction{})

	log.Info("database connection established")
	return db
}
