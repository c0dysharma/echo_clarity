package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/c0dysharma/echo_clarity/structs"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/calendar/v3"
)

var calendarService = structs.CalendarEvent{}

func GetCalendarEvents(c echo.Context) error {
	accessToken := ""

	// Fetch Google Calendar events for today
	events, err := calendarService.GetTodayEvents(accessToken)
	if err != nil {
		log.Printf("Error fetching calendar events: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch calendar events"})
	}

	// Return the events as a JSON response
	return c.JSON(http.StatusOK, events)
}

func CreateCalendarEvent(c echo.Context) error {
	accessToken := ""

	// Parse the request body
	var req structs.CalendarEvent
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate input
	if req.EventName == "" || req.StartTime == "" || req.EndTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	// Parse start and end times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("Invalid start time format: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start time format"})
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		log.Printf("Invalid end time format: %v", err)
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
	createdEvent, err := calendarService.CreateEvent(event, accessToken)
	if err != nil {
		log.Printf("Failed to create calendar event: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create event"})
	}

	// Return the created event as a response
	return c.JSON(http.StatusOK, createdEvent)
}
