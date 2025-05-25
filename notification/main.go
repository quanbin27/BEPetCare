package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/quanbin27/commons/config"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func startKafkaConsumer(service *Service, kafkaAddr string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaAddr},
		Topic:    "order_topic",
		GroupID:  "notification-service",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("Kafka consumer started, listening to order_topic...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading Kafka message: %v", err)
			continue
		}

		var order OrderData
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshaling Kafka message: %v", err)
			continue
		}
		log.Printf("Order received: %v", order)
		_, err = service.SendOrderConfirmationEmail(context.Background(), order.Email, order.OrderID, order.Items)
		if err != nil {
			log.Printf("Failed to send order confirmation email: %v", err)
		} else {
			log.Printf("Order confirmation email sent for order %d", order.OrderID)
		}
	}
}

func main() {
	dsn := config.Envs.NotificationsDSN
	log.Println("Connecting to database ...", dsn)
	kafkaAddr := config.Envs.KafkaAddr

	grpcAddr := config.Envs.Notification_Addr
	db, err := NewMySQLStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	db.AutoMigrate(&EmailNotification{})

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("failed to listen")
	}
	defer l.Close()

	store := NewMySQLNotificationStore(db)
	mailDialer := gomail.NewDialer("smtp.gmail.com", 587, config.Envs.EmailAddr, config.Envs.EmailPassword)
	service := NewService(store, mailDialer)
	NewGRPCHandler(grpcServer, service)

	go startKafkaConsumer(service, kafkaAddr)

	log.Println("Notification Service Listening on", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
