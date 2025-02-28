package main

import (
	"github.com/labstack/echo/v4"
	config "github.com/quanbin27/commons/config"
	orders "github.com/quanbin27/commons/genproto/orders"
	users "github.com/quanbin27/commons/genproto/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	httpAddr := config.Envs.HTTP_ADDR
	ordersServiceAddr := config.Envs.OrdersGrpcAddr
	ordersConn, err := grpc.NewClient(ordersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial order server: %v", err)
	}
	defer ordersConn.Close()
	usersServiceAddr := config.Envs.UsersGrpcAddr
	usersConn, err := grpc.NewClient(usersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}
	defer usersConn.Close()
	log.Println("Dialing orders service at ", ordersServiceAddr)
	orderClient := orders.NewOrderServiceClient(ordersConn)
	userClient := users.NewUserServiceClient(usersConn)
	e := echo.New()
	subrouter := e.Group("/api/v1")
	orderHandler := NewOrderHandler(orderClient)
	orderHandler.registerRoutes(subrouter)
	userHandler := NewUserHandler(userClient)
	userHandler.registerRoutes(subrouter)
	log.Println("Starting server on", httpAddr)
	if err := e.Start(httpAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
