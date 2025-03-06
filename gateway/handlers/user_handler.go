package handlers

import (
	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/users"
	"net/http"
)

type userHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(client pb.UserServiceClient) *userHandler {
	return &userHandler{client}
}
func (h *userHandler) registerRoutes(e *echo.Group) {
	e.GET("/hello", h.sayHello)
}
func (h *userHandler) sayHello(c echo.Context) error {

	return c.JSON(http.StatusOK, "1")
}
