package main

import (
	"github.com/labstack/echo/v4"
	"log"
)

const httpAddr = ":8080"

func main() {
	e := echo.New()
	subrouter := e.Group("/api/v1")
	httpHandler := NewHandler()
	httpHandler.registerRoutes(subrouter)
	log.Println("Starting server on", httpAddr)
	if err := e.Start(httpAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
