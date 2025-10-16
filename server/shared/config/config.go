// TODO: Implement shared configuration
package config

import (
	"os"
	"strconv"
)

type Config struct {
	PublicHost                      string
	Port                            string
	DatabaseURL                     string
	Database                        DatabaseConfig
	JWTSecret                       string
	JWTExpirationInSeconds          int64
	RefreshTokenExpirationInSeconds int64
	ResendAPIKey                    string
	FrontendURL                     string
	MobileVerificationURL           string
	OAuthConfig                     OAuthConfig
}

type DatabaseConfig struct {
	MaxConnections    int64
	MinConnections    int64
	MaxConnLifetime   int64 // in minutes
	MaxConnIdleTime   int64 // in minutes
	HealthCheckPeriod int64 // in minutes
	ConnectTimeout    int64 // in seconds
}

type OAuthConfig struct {
	GoogleClientID           string
	GoogleClientSecret       string
	GoogleRedirectURI        string
	GoogleMobileClientID     string
	GoogleMobileClientSecret string
	GoogleMobileRedirectURI  string

	GitHubClientID           string
	GitHubClientSecret       string
	GitHubRedirectURI        string
	GitHubMobileClientID     string
	GitHubMobileClientSecret string
	GitHubMobileRedirectURI  string

	FacebookClientID           string
	FacebookClientSecret       string
	FacebookRedirectURI        string
	FacebookMobileClientID     string
	FacebookMobileClientSecret string
	FacebookMobileRedirectURI  string

	OAuthStateSecret string
}

func NewConfig() Config {
	return LoadConfig()
}

func LoadConfig() Config {
	return Config{
		PublicHost:  getEnv("PUBLIC_HOST", "http://localhost"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Database: DatabaseConfig{
			MaxConnections:    getEnvAsInt("DB_MAX_CONNECTIONS", 30),
			MinConnections:    getEnvAsInt("DB_MIN_CONNECTIONS", 5),
			MaxConnLifetime:   getEnvAsInt("DB_MAX_CONN_LIFETIME", 60),  // 60 minutes
			MaxConnIdleTime:   getEnvAsInt("DB_MAX_CONN_IDLE_TIME", 30), // 30 minutes
			HealthCheckPeriod: getEnvAsInt("DB_HEALTH_CHECK_PERIOD", 5), // 5 minutes
			ConnectTimeout:    getEnvAsInt("DB_CONNECT_TIMEOUT", 10),    // 10 seconds
		},
		JWTExpirationInSeconds:          getEnvAsInt("JWT_EXP", 3600*24*7),
		RefreshTokenExpirationInSeconds: getEnvAsInt("REFRESH_TOKEN_EXP", 3600*24*30), // 30 days
		JWTSecret:                       getEnv("JWT_SECRET", ""),
		ResendAPIKey:                    getEnv("RESEND_API_KEY", ""),
		FrontendURL:                     getEnv("FRONTEND_URL", ""),
		MobileVerificationURL:           getEnv("MOBILE_VERIFICATION_URL", ""),
		OAuthConfig: OAuthConfig{
			GoogleClientID:             getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret:         getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURI:          getEnv("GOOGLE_REDIRECT_URI", ""),
			GoogleMobileClientID:       getEnv("GOOGLE_MOBILE_CLIENT_ID", ""),
			GoogleMobileClientSecret:   getEnv("GOOGLE_MOBILE_CLIENT_SECRET", ""),
			GoogleMobileRedirectURI:    getEnv("GOOGLE_MOBILE_REDIRECT_URI", ""),
			GitHubClientID:             getEnv("GITHUB_CLIENT_ID", ""),
			GitHubClientSecret:         getEnv("GITHUB_CLIENT_SECRET", ""),
			GitHubRedirectURI:          getEnv("GITHUB_REDIRECT_URI", ""),
			GitHubMobileClientID:       getEnv("GITHUB_MOBILE_CLIENT_ID", ""),
			GitHubMobileClientSecret:   getEnv("GITHUB_MOBILE_CLIENT_SECRET", ""),
			GitHubMobileRedirectURI:    getEnv("GITHUB_MOBILE_REDIRECT_URI", ""),
			FacebookClientID:           getEnv("FACEBOOK_CLIENT_ID", ""),
			FacebookClientSecret:       getEnv("FACEBOOK_CLIENT_SECRET", ""),
			FacebookRedirectURI:        getEnv("FACEBOOK_REDIRECT_URI", ""),
			FacebookMobileClientID:     getEnv("FACEBOOK_MOBILE_CLIENT_ID", ""),
			FacebookMobileClientSecret: getEnv("FACEBOOK_MOBILE_CLIENT_SECRET", ""),
			FacebookMobileRedirectURI:  getEnv("FACEBOOK_MOBILE_REDIRECT_URI", ""),
			OAuthStateSecret:           getEnv("OAUTH_STATE_SECRET", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}
	return fallback
}
