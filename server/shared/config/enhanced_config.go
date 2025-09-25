package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)



type ServiceConfig struct {
	Name        string
	Port        string
	Environment string
	LogLevel    string
	Debug       bool
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret         string
	ExpiresIn      time.Duration
	RefreshExpires time.Duration
}

type EnhancedConfig struct {
	Service        ServiceConfig
	Database       DatabaseConfig
	Redis          RedisConfig
	JWT            JWTConfig
	AuthServiceURL string
	UserServiceURL string
	AIServiceURL   string
	GatewayURL     string
}

func LoadEnhancedConfig() *EnhancedConfig {
	cfg := &EnhancedConfig{
		Service: ServiceConfig{
			Name:        getEnvStr("SERVICE_NAME", "lornian-service"),
			Port:        getEnvStr("PORT", "8080"),
			Environment: getEnvStr("ENVIRONMENT", "development"),
			LogLevel:    getEnvStr("LOG_LEVEL", "info"),
			Debug:       getEnvBoolExt("DEBUG", false),
		},
		Database: DatabaseConfig{
			MaxConnections:    getEnvAsInt("DB_MAX_CONNECTIONS", 30),
			MinConnections:    getEnvAsInt("DB_MIN_CONNECTIONS", 5),
			MaxConnLifetime:   getEnvAsInt("DB_MAX_CONN_LIFETIME", 60),
			MaxConnIdleTime:   getEnvAsInt("DB_MAX_CONN_IDLE_TIME", 30),
			HealthCheckPeriod: getEnvAsInt("DB_HEALTH_CHECK_PERIOD", 5),
			ConnectTimeout:    getEnvAsInt("DB_CONNECT_TIMEOUT", 10),
		},
		Redis: RedisConfig{
			URL:      getEnvStr("REDIS_URL", "redis://localhost:6379"),
			Password: getEnvStr("REDIS_PASSWORD", ""),
			DB:       int(getEnvAsInt("REDIS_DB", 0)),
		},
		JWT: JWTConfig{
			Secret:         getEnvStr("JWT_SECRET", "your-jwt-secret"),
			ExpiresIn:      getEnvDuration("JWT_EXPIRES_IN", 15*time.Minute),
			RefreshExpires: getEnvDuration("JWT_REFRESH_EXPIRES", 7*24*time.Hour),
		},
		AuthServiceURL: getEnvStr("AUTH_SERVICE_URL", "http://localhost:8001"),
		UserServiceURL: getEnvStr("USER_SERVICE_URL", "http://localhost:8002"),
		AIServiceURL:   getEnvStr("AI_SERVICE_URL", "http://localhost:8003"),
		GatewayURL:     getEnvStr("GATEWAY_URL", "http://localhost:8000"),
	}

	return cfg
}

func (c *EnhancedConfig) Validate() error {
	if c.JWT.Secret == "" || c.JWT.Secret == "your-jwt-secret" {
		return fmt.Errorf("JWT_SECRET must be set and not be the default value")
	}

	return nil
}

func getEnvStr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolExt(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
