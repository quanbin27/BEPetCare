package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/quanbin27/commons/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net"
)

func NewMySQLStorage(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
func initStorage(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")
}
func main() {
	dsn := config.Envs.UsersDSN
	log.Println("Connecting to database ...", dsn)
	grpcAddr := config.Envs.UsersGrpcAddr
	db, err := NewMySQLStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	db.AutoMigrate(User{}, Role{}, UserRole{}, EmployeeBranch{})
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("failed to listen")
	}
	defer l.Close()
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	notificationsConn, err := grpc.NewClient(config.Envs.Notification_Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial order server: %v", err)
	}
	defer notificationsConn.Close()
	userStore := NewStore(db)
	userService := NewService(userStore, redisClient, notificationsConn, config.Envs.BaseURL)
	NewGrpcUsersHandler(grpcServer, userService)
	log.Println("User Service Listening on", grpcAddr)
	grpcServer.Serve(l)
}
