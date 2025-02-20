package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}
func (h *handler) registerRoutes(e *echo.Group) {
	e.GET("/hello", h.sayHello)
}
func (h *handler) sayHello(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello world")
}
