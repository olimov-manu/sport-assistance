package configs

import (
	"fmt"
	"log"
	"os"
	"sport-assistance/pkg/utils"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	DBMaxConn          int
	DBConnectionString string
}

type ServerConfig struct {
	Port         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type SecurityConfig struct {
	JWTSecretKey       string
	AccessTokenTTL     time.Duration
	AccessTokenSecret  string
	RefreshTokenTTL    time.Duration
	RefreshTokenSecret string
}

type LoggerConfig struct {
	Level string
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DBName   string
}

type Config struct {
	ServerConfig   ServerConfig
	DatabaseConfig DatabaseConfig
	SecurityConfig SecurityConfig
	Logger         LoggerConfig
	RedisConfig    RedisConfig
}

// GetConfigs загружает конфигурацию из .env или переменных окружения
func GetConfigs() *Config {
	// Загружаем .env, если он есть
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		port = 5432
	}

	maxConn, err := strconv.Atoi(getEnv("DB_MAX_CONN", "10"))
	if err != nil {
		maxConn = 10
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		port,
		getEnv("DB_NAME", "postgres"),
		getEnv("DB_SSLMODE", "disable"),
	)

	return &Config{
		ServerConfig: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			WriteTimeout: utils.ToDuration(getEnv("WRITE_TIMEOUT", "30s")),
			ReadTimeout:  utils.ToDuration(getEnv("READ_TIMEOUT", "30s")),
		},
		DatabaseConfig: DatabaseConfig{
			DBHost:             getEnv("DB_HOST", "localhost"),
			DBPort:             port,
			DBUser:             getEnv("DB_USER", "postgres"),
			DBPassword:         getEnv("DB_PASSWORD", ""),
			DBName:             getEnv("DB_NAME", "postgres"),
			DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
			DBMaxConn:          maxConn,
			DBConnectionString: dsn,
		},
		RedisConfig: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DBName:   getEnv("REDIS_DB", ""),
		},
		SecurityConfig: SecurityConfig{
			JWTSecretKey:       getEnv("JWT_SECRET_KEY", ""),
			AccessTokenTTL:     utils.ToDuration(getEnv("ACCESS_TOKEN_TTL", "10m")),
			AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", ""),
			RefreshTokenTTL:    utils.ToDuration(getEnv("REFRESH_TOKEN_TTL", "720h")),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", ""),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv возвращает значение переменной окружения или дефолтное значение
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
