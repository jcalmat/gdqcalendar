package main

import (
	"context"
	"log"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	parserapp "github.com/jcalmat/gdqcalendar/pkg/parser/app"
)

func main() {
	ctx := context.Background()

	// setup external services
	gconfig, err := parseGoogleConfig()
	if err != nil {
		log.Fatalf("failed to parse google config: %s", err.Error())
	}
	client := getClient(gconfig)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// parse config
	config := parseConfig()

	c := &Gcal{
		Service:    srv,
		CalendarID: config.CalendarID,
		GDQParser: parserapp.App{
			ScheduleURL:    "https://gamesdonequick.com/schedule/",
			ScheduleApiURL: "https://gamesdonequick.com/tracker/api/v2/events/{version}/runs?",
		},
	}

	err = c.reimportEvents()
	if err != nil {
		log.Fatalf("failed to reimport events: %s", err.Error())
	}
}
