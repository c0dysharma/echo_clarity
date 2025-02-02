// File: handlers/auth.go
package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/c0dysharma/echo_clarity/helpers"
	"github.com/c0dysharma/echo_clarity/models"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

// GoogleLoginHandler redirects to Google's OAuth login
func GoogleLoginHandler(c echo.Context) error {
	gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	return nil
}

func GoogleCallbackHandler(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// encrypt user refresh token
	hashedAccessToken, err1 := helpers.EncryptPassword(user.AccessToken)
	hashedRefreshToken, err2 := helpers.EncryptPassword(user.RefreshToken)
	if err1 != nil || err2 != nil {
		log.Error("error in password hash")
	}

	// store in db if not exists
	var dbuser models.User
	helpers.DB.Where("email = ?", user.Email).First(&dbuser)

	if dbuser.Email == "" {
		// create new user
		dbuser = models.User{
			Email:                user.Email,
			Name:                 user.Name,
			AccessToken:          hashedAccessToken,
			AccessTokenExpiresAt: user.ExpiresAt,
		}
		if user.RefreshToken != "" {
			dbuser.RefreshToken = hashedRefreshToken
		}

		helpers.DB.Create(&dbuser)
	} else {
		// else update refresh token
		dbuser.AccessToken = hashedAccessToken
		dbuser.AccessTokenExpiresAt = user.ExpiresAt

		if user.RefreshToken != "" {
			dbuser.RefreshToken = hashedRefreshToken
		}
		helpers.DB.Save(&dbuser)
	}

	// Generate JWT token
	token, err := helpers.GenerateJWT(dbuser.ID, dbuser.Email, dbuser.JWTTokenVersion)
	if err!= nil {
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
  }

	// update the token version for user
	dbuser.JWTTokenVersion++
  helpers.DB.Save(&dbuser)

  // Return both user and token
  return c.JSON(http.StatusOK, map[string]interface{}{
    "user":  dbuser,
    "token": token,
    "gUser": user,
  })
}
