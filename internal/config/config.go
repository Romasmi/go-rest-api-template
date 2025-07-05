package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxConnections  int
	MinConnections  int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type JWTConfig struct {
	Secret        string
	ExpirationTTL time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	readTimeout := getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second)
	writeTimeout := getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second)
	idleTimeout := getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second)

	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/go_rest_api?sslmode=disable")
	maxConnections := getEnvAsInt("DATABASE_MAX_CONNECTIONS", 10)
	minConnections := getEnvAsInt("DATABASE_MIN_CONNECTIONS", 2)
	maxConnLifetime := getEnvAsDuration("DATABASE_MAX_CONN_LIFETIME", time.Hour)
	maxConnIdleTime := getEnvAsDuration("DATABASE_MAX_CONN_IDLE_TIME", 30*time.Minute)

	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	jwtExpirationTTL := getEnvAsDuration("JWT_EXPIRATION_TTL", 24*time.Hour)

	config := &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		Database: DatabaseConfig{
			URL:             dbURL,
			MaxConnections:  maxConnections,
			MinConnections:  minConnections,
			MaxConnLifetime: maxConnLifetime,
			MaxConnIdleTime: maxConnIdleTime,
		},
		JWT: JWTConfig{
			Secret:        jwtSecret,
			ExpirationTTL: jwtExpirationTTL,
		},
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
