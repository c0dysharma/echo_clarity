// File: handlers/auth.go
package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// GoogleLoginHandler redirects to Google's OAuth login
func GoogleLoginHandler(c echo.Context) error {
	fmt.Println(os.Getenv("GOOGLE_CLIENT_KEY"),
	os.Getenv("GOOGLE_CLIENT_SECRET"),
	os.Getenv("GOOGLE_REDIRECT_URL"))

	gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	return nil
}

// GoogleCallbackHandler handles the callback from Google's OAuth
func GoogleCallbackHandler(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}
