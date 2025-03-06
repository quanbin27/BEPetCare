package handlers

import (
	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/orders"
	"net/http"
)

type orderHandler struct {
	client pb.OrderServiceClient
}

func NewOrderHandler(client pb.OrderServiceClient) *orderHandler {
	return &orderHandler{client}
}
func (h *orderHandler) registerRoutes(e *echo.Group) {
	e.GET("/hello", h.sayHello)
}
func (h *orderHandler) sayHello(c echo.Context) error {

	return c.JSON(http.StatusOK, "hello")
}
