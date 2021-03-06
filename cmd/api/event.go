package main

import (
	"bufio"
	"fmt"
	"os"
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

	ids := make([]string, 0)

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
				DateTime: game.StartDate.Add(game.Duration).Format(time.RFC3339),
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
			Description: fmt.Sprintf("Runners: %s\nRun Time: %s\nCategory: %s\nHost: %s\nSetup Time: %s", strings.Join(game.Runners, ", "), fmtDuration(game.Duration), game.Category, game.Host, fmtDuration(game.SetupDuration)),
		}

		err := retry.Do(func() error {
			ret, err := c.Service.Events.Insert(c.CalendarID, event).Do()
			if err != nil {
				fmt.Printf("Error creating event: %s\n", err)
				fmt.Println("Retry in 10 seconds...")
				return err
			}
			ids = append(ids, ret.Id)
			return nil
		}, retry.Attempts(3), retry.Delay(time.Second*10))
		if err != nil {
			return err
		}
	}

	err = writeToFile(ids)
	if err != nil {
		return err
	}

	return nil
}

// writeToFile writes the ids to a file in order to delete them later
func writeToFile(ids []string) error {
	var f *os.File
	f, err := os.OpenFile("ids.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintln(strings.Join(ids, "\n")))
	if err != nil {
		return err
	}

	return nil
}

// deleteEvents deletes all events in the calendar based on the ids in the ids.txt file
func (c *Gcal) deleteEvents() error {
	f, err := os.Open("ids.txt")
	if err != nil {
		return nil
	}

	defer f.Close()

	ids := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, id := range ids {
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

	// clean up the file
	err = os.Truncate("ids.txt", 0)
	if err != nil {
		return err
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
