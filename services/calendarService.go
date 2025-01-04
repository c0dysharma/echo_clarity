package services

import (
	"context"
	"log"
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

func (c *CalendarEvent) GetTodayEvents(accessToken string) ([]*calendar.Event, error) {
	// Create a new Calendar Service using the access token
	srv, err := c.getCalendarService(accessToken)
	if err != nil {
		log.Printf("Unable to create calendar service: %v", err)
		return nil, err
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
		log.Printf("Unable to retrieve calendar events: %v", err)
		return nil, err
	}

	return events.Items, nil
}

func (c *CalendarEvent) CreateEvent(event *calendar.Event, accessToken string) (*calendar.Event, error) {
	srv, err := c.getCalendarService(accessToken)
	if err != nil {
		log.Printf("Unable to create calendar service: %v", err)
		return nil, err
	}

	// Insert the event into the user's primary calendar
	createdEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		log.Printf("Unable to create event: %v", err)
		return nil, err
	}

	return createdEvent, nil
}