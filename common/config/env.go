package config

import (
	"github.com/lpernett/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DSN                    string
	JWTExpirationInSeconds int64
	JWTSecret              string
	HTTP_ADDR              string
	UserGrpcAddr           string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Không thể load file .env, sử dụng biến môi trường hệ thống")
	}
	return Config{
		DSN:                    getEnv("DSN", ""),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not-secret-anymore?"),
		HTTP_ADDR:              getEnv("HTTP_ADDR", ":8080"),
		UserGrpcAddr:           getEnv("USER_GRPC_ADDR", ":8081"),
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
