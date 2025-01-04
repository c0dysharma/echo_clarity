// File: cmd/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/c0dysharma/echo_clarity/handlers"
	"github.com/c0dysharma/echo_clarity/helpers"
	"github.com/c0dysharma/echo_clarity/middlewares"
	"github.com/c0dysharma/echo_clarity/models"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func initOAuth() {

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_KEY"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_REDIRECT_URL"),
			"email", "profile", "https://www.googleapis.com/auth/calendar",
		),
	)

	// Override the default GetProviderName function
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "google", nil
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	store := sessions.NewCookieStore([]byte(os.Getenv("GOOGLE_CLIENT_KEY")))
	store.MaxAge(86400 * 30) // 30 days
	gothic.Store = store

	// initialize OAuth
	initOAuth()

	// initialize database
	helpers.ConnectDB()

	// auto migrate models
	db := helpers.DB
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate models")
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", handlers.PongHandler)
	e.GET("/auth/google", handlers.GoogleLoginHandler)
	e.GET("/auth/google/callback", handlers.GoogleCallbackHandler)

	// Authenticated routes
	authGroup := e.Group("")
	authGroup.Use(middlewares.AuthMiddleware)
	authGroup.GET("/calendar", handlers.GetCalendarEvents)
	authGroup.POST("/calendar", handlers.CreateCalendarEvent)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}
