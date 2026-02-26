package common

import (
	"os"
	"strconv"
)

var (
	Port                = getEnv("PORT", "3000")
	ServerURL           = getEnv("SERVER_URL", "http://localhost:3000")
	SessionSecret       = getEnv("SESSION_SECRET", "default-session-secret")
	JWTSecret           = getEnv("JWT_SECRET", "default-jwt-secret")
	GinMode             = getEnv("GIN_MODE", "debug")
	SqlDSN              = getEnv("SQL_DSN", "")
	CPAUpstreamURL      = getEnv("CPA_UPSTREAM_URL", "")
	CPAUpstreamKey      = getEnv("CPA_UPSTREAM_KEY", "")
	LinuxDOClientID     = getEnv("LINUXDO_CLIENT_ID", "")
	LinuxDOClientSecret = getEnv("LINUXDO_CLIENT_SECRET", "")
	CORSAllowOrigins    = getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173")
	TrustedProxies      = getEnv("TRUSTED_PROXIES", "")
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
