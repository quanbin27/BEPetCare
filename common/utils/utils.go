package utils

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func GetTokenFromRequest(c echo.Context) string {
	tokenAuth := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(tokenAuth, "Bearer ") {
		return strings.TrimPrefix(tokenAuth, "Bearer ")
	}

	tokenQuery := c.QueryParam("token")
	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}
