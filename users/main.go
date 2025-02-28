package main

import (
	"github.com/quanbin27/commons/config"
	"google.golang.org/grpc"
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
	userStore := NewStore(db)
	userService := NewService(userStore)
	NewGrpcUsersHandler(grpcServer, userService)
	log.Println("User Service Listening on", grpcAddr)
	grpcServer.Serve(l)
}
