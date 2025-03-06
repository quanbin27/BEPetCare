package main

import (
	"github.com/labstack/echo/v4"
	"github.com/quanbin27/BEPetCare-gateway/handlers"
	config "github.com/quanbin27/commons/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	httpAddr := config.Envs.HTTP_ADDR
	e := echo.New()
	subrouter := e.Group("/api/v1")
	usersServiceAddr := config.Envs.UsersGrpcAddr
	ordersServiceAddr := config.Envs.OrdersGrpcAddr
	recordsServiceAddr := config.Envs.RecordsGrpcAddr
	productsServiceAddr := config.Envs.ProductsGrpcAddr
	paymentsServiceAddr := config.Envs.PaymentsGrpcAddr
	appointmentsServiceAddr := config.Envs.AppointmentsGrpcAddr
	petRecordConn, err := grpc.NewClient(recordsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	usersConn, err := grpc.NewClient(usersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	productsConn, err := grpc.NewClient(productsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	paymentsConn, err := grpc.NewClient(paymentsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	ordersConn, err := grpc.NewClient(ordersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial order server: %v", err)
	}

	appointmentsConn, err := grpc.NewClient(appointmentsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial order server: %v", err)
	}
	orderHandler := handlers.NewOrderHandler(orderClient)
	orderHandler.registerRoutes(subrouter)
	userHandler := handlers.NewUserHandler(userClient)
	userHandler.registerRoutes(subrouter)
	log.Println("Starting server on", httpAddr)
	if err := e.Start(httpAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
