package main

import (
	"github.com/payOSHQ/payos-lib-golang"
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
	err := payos.Key(config.Envs.ClientID, config.Envs.ApiKey, config.Envs.CheckSumKey)
	if err != nil {
		log.Fatal("Failed to initialize PayOS keys:", err)
	}
	dsn := config.Envs.PaymentsDSN
	log.Println("Connecting to database ...", dsn)
	grpcAddr := config.Envs.PaymentsGrpcAddr
	db, err := NewMySQLStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	db.AutoMigrate(Payment{})
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("failed to listen")
	}
	defer l.Close()
	paymentStore := NewStore(db)
	paymentService := NewPaymentService(paymentStore)
	NewPaymentGrpcHandler(grpcServer, paymentService)
	log.Println("Payments Service Listening on", grpcAddr)
	grpcServer.Serve(l)
}
