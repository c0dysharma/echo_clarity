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

	// Fix log statements with proper key-value pairs
	log.Info("token info", "access_token", rUser.AccessToken)
	log.Info("token info", "refresh_token", rUser.RefreshToken)

	// Fetch Google Calendar events using decrypted token
	events, err, isOperational := calendarService.GetTodayEvents(rUser.AccessToken)
	if err != nil {
		// Fix error logging with proper key-value pairs
		log.Error("calendar event fetch failed", "error", err)

		eC := http.StatusUnauthorized
		if isOperational {
			eC = http.StatusInternalServerError
		}

		return c.JSON(eC, map[string]string{"error": "Failed to fetch calendar events"})
	}

	return c.JSON(http.StatusOK, events)
}

func CreateCalendarEvent(c echo.Context) error {
	rUser := c.Get("user").(models.User)

	if rUser.Email == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	log.Info("token info", "access_token", rUser.AccessToken)
	log.Info("token info", "refresh_token", rUser.RefreshToken)

	// Parse the request body
	var req structs.CalendarEvent
	if err := c.Bind(&req); err != nil {
		log.Error("Failed to parse request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate input
	if req.EventName == "" || req.StartTime.IsZero() || req.EndTime.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	// Create a new event
	event := &calendar.Event{
		Summary: req.EventName,
		Start: &calendar.EventDateTime{
			DateTime: req.StartTime.Format(time.RFC3339),
			TimeZone: req.StartTime.Location().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: req.EndTime.Format(time.RFC3339),
			TimeZone: req.EndTime.Location().String(),
		},
	}

	// Call the helper function to create the event
	createdEvent, err, isOperational := calendarService.CreateEvent(event, rUser.AccessToken)
	if err != nil {
		log.Error("Failed to create calendar event", "error", err)

		eC := http.StatusUnauthorized
		if isOperational {
			eC = http.StatusInternalServerError
		}

		return c.JSON(eC, map[string]string{"error": "Failed to create event"})
	}

	// Return the created event as a response
	return c.JSON(http.StatusOK, createdEvent)
}
