package config

import (
	"github.com/lpernett/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	UsersDSN               string
	ProductsDSN            string
	OrdersDSN              string
	AppointmentsDSN        string
	RecordsDSN             string
	PaymentsDSN            string
	JWTExpirationInSeconds int64
	JWTSecret              string
	HTTP_ADDR              string
	UsersGrpcAddr          string
	ProductsGrpcAddr       string
	OrdersGrpcAddr         string
	AppointmentsGrpcAddr   string
	RecordsGrpcAddr        string
	PaymentsGrpcAddr       string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Không thể load file .env, sử dụng biến môi trường hệ thống")
	}
	return Config{
		UsersDSN:        getEnv("USERS_DSN", ""),
		ProductsDSN:     getEnv("PRODUCTS_DSN", ""),
		OrdersDSN:       getEnv("ORDERS_DSN", ""),
		AppointmentsDSN: getEnv("APPOINTMENTS_DSN", ""),
		RecordsDSN:      getEnv("RECORDS_DSN", ""),
		PaymentsDSN:     getEnv("PAYMENTS_DSN", ""),

		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not-secret-anymore?"),
		HTTP_ADDR:              getEnv("HTTP_ADDR", ":8080"),

		UsersGrpcAddr:        getEnv("USER_GRPC_ADDR", ":8081"),
		ProductsGrpcAddr:     getEnv("PRODUCTS_GRPC_ADDR", ":8082"),
		OrdersGrpcAddr:       getEnv("ORDERS_GRPC_ADDR", ":8083"),
		AppointmentsGrpcAddr: getEnv("APPOINTMENTS_GRPC_ADDR", ":8084"),
		RecordsGrpcAddr:      getEnv("RECORDS_GRPC_ADDR", ":8085"),
		PaymentsGrpcAddr:     getEnv("PAYMENTS_GRPC_ADDR", ":8086"),
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
