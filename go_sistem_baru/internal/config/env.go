package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Web      WebConfig
	Database DatabaseConfig
	Redis    RedisConfig
	App      AppConfig
}

type WebConfig struct {
	Port        int
	Prefork     bool
	MetricsPort int
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Name            string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Addr              string
	Password          string
	DB                int
	PoolSize          int
	InquiryTTLSeconds int
}

type AppConfig struct {
	LogLevel  string
	JwtSecret string
}

func LoadConfig() *Config {
	cfg := &Config{
		Web: WebConfig{
			Port:        3000,
			Prefork:     false,
			MetricsPort: 8080,
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			Username:        "qris_user",
			Password:        "qris_pass",
			Name:            "qris_db",
			SSLMode:         "disable",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 3600,
		},
		Redis: RedisConfig{
			Addr:              "localhost:6379",
			Password:          "",
			DB:                0,
			PoolSize:          100,
			InquiryTTLSeconds: 300,
		},
		App: AppConfig{
			LogLevel:  "info",
			JwtSecret: "supersecretkey",
		},
	}

	cfg.Web.Port = getEnvInt("WEB_PORT", cfg.Web.Port)
	cfg.Web.Prefork = getEnvBool("WEB_PREFORK", cfg.Web.Prefork)
	cfg.Web.MetricsPort = getEnvInt("METRICS_PORT", cfg.Web.MetricsPort)

	cfg.Database.Host = getEnvString("DATABASE_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvInt("DATABASE_PORT", cfg.Database.Port)
	cfg.Database.Username = getEnvString("DATABASE_USERNAME", cfg.Database.Username)
	cfg.Database.Password = getEnvString("DATABASE_PASSWORD", cfg.Database.Password)
	cfg.Database.Name = getEnvString("DATABASE_NAME", cfg.Database.Name)
	cfg.Database.SSLMode = getEnvString("DATABASE_SSLMODE", cfg.Database.SSLMode)
	cfg.Database.MaxIdleConns = getEnvInt("DATABASE_MAX_IDLE_CONNS", cfg.Database.MaxIdleConns)
	cfg.Database.MaxOpenConns = getEnvInt("DATABASE_MAX_OPEN_CONNS", cfg.Database.MaxOpenConns)
	cfg.Database.ConnMaxLifetime = getEnvInt("DATABASE_CONN_MAX_LIFETIME", cfg.Database.ConnMaxLifetime)

	cfg.Redis.Addr = getEnvString("REDIS_ADDR", cfg.Redis.Addr)
	cfg.Redis.Password = getEnvString("REDIS_PASSWORD", cfg.Redis.Password)
	cfg.Redis.DB = getEnvInt("REDIS_DB", cfg.Redis.DB)
	cfg.Redis.PoolSize = getEnvInt("REDIS_POOL_SIZE", cfg.Redis.PoolSize)
	cfg.Redis.InquiryTTLSeconds = getEnvInt("REDIS_INQUIRY_TTL_SECONDS", cfg.Redis.InquiryTTLSeconds)

	cfg.App.LogLevel = getEnvString("APP_LOG_LEVEL", cfg.App.LogLevel)
	cfg.App.JwtSecret = getEnvString("APP_JWT_SECRET", cfg.App.JwtSecret)

	return cfg
}

func getEnvString(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.Atoi(strings.TrimSpace(val))
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "1", "true", "yes", "y":
			return true
		case "0", "false", "no", "n":
			return false
		}
	}
	return fallback
}
