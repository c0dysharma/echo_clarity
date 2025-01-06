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
		tokenVersion := float32(claims["tokenVersion"].(float64))
		var user models.User

		helpers.DB.Find(&user, userID)
		if user.ID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}

		if user.JWTTokenVersion != tokenVersion{
			return echo.NewHTTPError(http.StatusUnauthorized, "token version expired")
		}

		c.Set("user", user)
		return next(c)
	}
}

func DecryptToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rUser := c.Get("user").(models.User)

		if rUser.Email == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
		}

		// decrypt the token and set user again

		dAT, err := helpers.DecryptPassword(rUser.AccessToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decrypt access token"})
		}

		dRT, err := helpers.DecryptPassword(rUser.RefreshToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decrypt refresh token"})
		}

		rUser.AccessToken = dAT
		rUser.RefreshToken = dRT

		c.Set("user", rUser)
		return next(c)
	}
}
