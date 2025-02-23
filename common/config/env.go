package config

import (
	"github.com/lpernett/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DSN                    string
	JWTExpirationInSeconds int64
	JWTSecret              string
	HTTP_ADDR              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		DSN:                    getEnv("DSN", ""),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not-secret-anymore?"),
		HTTP_ADDR:              getEnv("HTTP_ADDR", ":8080"),
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
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
		return fallback
	}
	return fallback
}
