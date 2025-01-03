// File: handlers/pong.go
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PongHandler is a simple ping-pong handler for health checks
func PongHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
