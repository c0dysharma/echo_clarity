package middlewares

import (
	"net/http"
	"strings"

	"github.com/c0dysharma/echo_clarity/helpers"
	"github.com/c0dysharma/echo_clarity/models"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
		}

		claims, err := helpers.VerifyToken(tokenParts[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		// Access ID from MapClaims
		userID := claims["ID"]
		var user models.User

		helpers.DB.Find(&user, userID)
		if user.ID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}

		c.Set("user", user)
		return next(c)
	}
}
