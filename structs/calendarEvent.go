package structs

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarEvent struct {
	EventName string `json:"event_name"`
	StartTime string `json:"start_time"` // ISO 8601 format: "2025-01-03T10:00:00-05:00"
	EndTime   string `json:"end_time"`   // ISO 8601 format: "2025-01-03T11:00:00-05:00"
}

func (c *CalendarEvent) getCalendarService(accessToken string) (*calendar.Service, error) {
	// Create a context
	ctx := context.Background()

	// Create a new Calendar Service using the access token
	srv, err := calendar.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)))

	return srv, err
}

/*
*
[]*calendar.Event - Calendar events
error - Error if any
bool - Boolean indicating whether its Operational Error or not
*/
func (c *CalendarEvent) GetTodayEvents(accessToken string) ([]*calendar.Event, error, bool) {
	// Create a new Calendar Service using the access token
	srv, err := c.getCalendarService(accessToken)
	if err != nil {
		return nil, err, true
	}

	// Get the current date in RFC3339 format
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location()).Format(time.RFC3339)

	// Fetch events from the calendar for today
	events, err := srv.Events.List("primary").
		TimeMin(startOfDay).
		TimeMax(endOfDay).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, err, false
	}

	return events.Items, nil, false
}

func (c *CalendarEvent) CreateEvent(event *calendar.Event, accessToken string) (*calendar.Event, error, bool) {
	srv, err := c.getCalendarService(accessToken)
	if err != nil {
		return nil, err, true
	}

	// Insert the event into the user's primary calendar
	createdEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, err, false
	}

	return createdEvent, nil, false
}
