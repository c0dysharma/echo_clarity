package handlers

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"

	"github.com/c0dysharma/echo_clarity/models"
	"github.com/c0dysharma/echo_clarity/structs"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/calendar/v3"
)

var calendarService = structs.CalendarEvent{}

func GetCalendarEvents(c echo.Context) error {
	// Get email from context
	rUser := c.Get("user").(models.User)

	if rUser.Email == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	log.Info("User AccessToken", rUser.AccessToken)
	log.Info("User RefreshToken", rUser.RefreshToken)

	// Fetch Google Calendar events using decrypted token
	events, err := calendarService.GetTodayEvents(rUser.AccessToken)
	if err != nil {
		log.Error("Error fetching calendar events: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch calendar events"})
	}

	return c.JSON(http.StatusOK, events)
}

func CreateCalendarEvent(c echo.Context) error {
	rUser := c.Get("user").(models.User)

	if rUser.Email == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	log.Info("User AccessToken", rUser.AccessToken)
	log.Info("User RefreshToken", rUser.RefreshToken)

	// Parse the request body
	var req structs.CalendarEvent
	if err := c.Bind(&req); err != nil {
		log.Error("Failed to parse request body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate input
	if req.EventName == "" || req.StartTime == "" || req.EndTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	// Parse start and end times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Error("Invalid start time format: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start time format"})
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		log.Error("Invalid end time format: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end time format"})
	}

	// Create a new event
	event := &calendar.Event{
		Summary: req.EventName,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: startTime.Location().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: endTime.Location().String(),
		},
	}

	// Call the helper function to create the event
	createdEvent, err := calendarService.CreateEvent(event, rUser.AccessToken)
	if err != nil {
		log.Error("Failed to create calendar event: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create event"})
	}

	// Return the created event as a response
	return c.JSON(http.StatusOK, createdEvent)
}
