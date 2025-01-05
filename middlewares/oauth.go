package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/c0dysharma/echo_clarity/helpers"
	"github.com/c0dysharma/echo_clarity/models"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

type OAuthRefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func getNewAccessToken(refreshToken string) (OAuthRefreshResponse, error) {

	reqBody := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     os.Getenv("GOOGLE_CLIENT_KEY"),
		"client_secret": os.Getenv("GOOGLE_CLIENT_SECRET"),
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return OAuthRefreshResponse{}, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("https://accounts.google.com/o/oauth2/token",
		"application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return OAuthRefreshResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var oauthResp OAuthRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return OAuthRefreshResponse{}, fmt.Errorf("error decoding response: %v", err)
	}

	return oauthResp, nil
}

func RefreshOAuthAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rUser := c.Get(("user")).(models.User)

		if rUser.AccessToken == "" || rUser.RefreshToken == "" || rUser.AccessTokenExpiresAt == (time.Time{}) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Re-login required"})
		}

		// check if access token's expiry is past
		if rUser.AccessTokenExpiresAt.Before(time.Now()) {
			// refresh token
			log.Info("oAuth token expired...Refreshing")

			response, err := getNewAccessToken(rUser.RefreshToken)

			log.Info(response)

			// if refreshing failed i.e refresh token is expired
			// return 401
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}

			// Calculate expiry time from expires_in seconds
			expiryTime := time.Now().Add(time.Second * time.Duration(response.ExpiresIn))

			// save in context
			rUser.AccessToken = response.AccessToken
			rUser.AccessTokenExpiresAt = expiryTime
			c.Set("user", rUser)

			// save in db
			var user models.User
			helpers.DB.Find(&user, rUser.ID)

			eP, err := helpers.EncryptPassword(rUser.AccessToken)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encrypt access token"})
			}

			helpers.DB.Model(&user).Updates(models.User{AccessToken: eP, AccessTokenExpiresAt: rUser.AccessTokenExpiresAt})

		}
		return next(c)
	}
}
