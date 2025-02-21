package main

import (
	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/api"
	"net/http"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}
func (h *handler) registerRoutes(e *echo.Group) {
	e.GET("/hello", h.sayHello)
}
func (h *handler) sayHello(c echo.Context) error {
	h.client.CreateOrder(c, &pb.CreateOrderRequest{CustomerID: "1"})
	return c.JSON(http.StatusOK, "hello world")
}
