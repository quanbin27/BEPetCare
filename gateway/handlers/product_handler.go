package handlers

import (
	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/products"
	"net/http"
)

type productHandler struct {
	client pb.ProductServiceClient
}

func NewProductHandler(client pb.ProductServiceClient) *productHandler {
	return &productHandler{client}
}
func (h *productHandler) registerRoutes(e *echo.Group) {
	e.GET("/hello", h.sayHello)
}
func (h *productHandler) sayHello(c echo.Context) error {

	return c.JSON(http.StatusOK, "")
}
