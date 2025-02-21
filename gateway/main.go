package main

import (
	"github.com/labstack/echo/v4"
	common "github.com/quanbin27/commons"
	pb "github.com/quanbin27/commons/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	httpAddr          = common.EnvString("HTTP_ADDR", ":8080")
	ordersServiceAddr = "localhost:3000"
)

func main() {
	conn, err := grpc.NewClient(ordersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()
	log.Println("Dialing orders service at ", ordersServiceAddr)
	c := pb.NewOrderServiceClient(conn)
	e := echo.New()
	subrouter := e.Group("/api/v1")
	httpHandler := NewHandler(c)
	httpHandler.registerRoutes(subrouter)
	log.Println("Starting server on", httpAddr)
	if err := e.Start(httpAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
