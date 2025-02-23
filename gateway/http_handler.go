package main

import (
	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/orders"
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
	order, err := h.client.CreateOrder(c.Request().Context(), &pb.CreateOrderRequest{CustomerID: "1"})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, order)
}
