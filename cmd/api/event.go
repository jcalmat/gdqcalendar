package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"google.golang.org/api/calendar/v3"
)

func (c *Gcal) reimportEvents() error {
	err := c.deleteEvents()
	if err != nil {
		return err
	}

	return c.createEvents()
}

// createEvents parses the schedule and creates events in the calendar
func (c *Gcal) createEvents() error {

	cal, err := c.GDQParser.Parse()
	if err != nil {
		return err
	}

	for i, game := range cal.Games {
		var end *calendar.EventDateTime
		if i < len(cal.Games)-1 {
			end = &calendar.EventDateTime{
				DateTime: cal.Games[i+1].StartDate.Format(time.RFC3339),
				TimeZone: "America/New_York",
			}
		} else {
			end = &calendar.EventDateTime{
				DateTime: game.EndDate.Format(time.RFC3339),
				TimeZone: "America/New_York",
			}
		}

		event := &calendar.Event{
			Summary: game.Name,
			Start: &calendar.EventDateTime{
				DateTime: game.StartDate.Format(time.RFC3339),
				TimeZone: "America/New_York",
			},
			End:         end,
			Description: fmt.Sprintf("Runners: %s\nRun Time: %s\nCategory: %s\nHost: %s\nSetup Time: %s", strings.Join(game.Runners, ", "), game.Duration, game.Category, game.Host, game.SetupDuration),
		}

		err := retry.Do(func() error {
			_, err := c.Service.Events.Insert(c.CalendarID, event).Do()
			if err != nil {
				fmt.Printf("Error creating event: %s\n", err)
				fmt.Println("Retry in 10 seconds...")
				return err
			}
			return nil
		}, retry.Attempts(3), retry.Delay(time.Second*10))
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteEvents deletes all events in the calendar based on the ids in the ids.txt file
func (c *Gcal) deleteEvents() error {
	// clean events happening in the future
	events, err := c.Service.Events.List(c.CalendarID).SingleEvents(true).TimeMin(time.Now().Format(time.RFC3339)).Do()
	if err != nil {
		return err
	}

	for _, ev := range events.Items {
		id := ev.Id
		err := retry.Do(func() error {
			err := c.Service.Events.Delete(c.CalendarID, id).Do()
			if err != nil {
				fmt.Printf("Error deleting event: %s\n", err)
				fmt.Println("Retry in 10 seconds...")
				return err
			}
			return nil
		}, retry.Attempts(3), retry.Delay(time.Second*10))
		if err != nil {
			return err
		}
	}

	return nil
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
