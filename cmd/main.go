// File: cmd/main.go
package main

import (
	"net/http"
	"os"

	"github.com/c0dysharma/echo_clarity/handlers"
	"github.com/c0dysharma/echo_clarity/helpers"
	"github.com/c0dysharma/echo_clarity/middlewares"
	"github.com/c0dysharma/echo_clarity/models"
	"github.com/charmbracelet/log"
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

	// Initialize logger
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05",
	})

	e := echo.New()
	e.HideBanner = false
	e.HidePort = false

	// Custom logging middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
			)
			return nil
		},
	}))

	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", handlers.PongHandler)
	e.GET("/auth/google", handlers.GoogleLoginHandler)
	e.GET("/auth/google/callback", handlers.GoogleCallbackHandler)

	// Authenticated routes
	authGroup := e.Group("")
	authGroup.Use(middlewares.AuthMiddleware)
	authGroup.Use(middlewares.DecryptToken)
	authGroup.Use(middlewares.RefreshOAuthAccessToken)
	authGroup.GET("/calendar", handlers.GetCalendarEvents)
	authGroup.POST("/calendar", handlers.CreateCalendarEvent)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}
