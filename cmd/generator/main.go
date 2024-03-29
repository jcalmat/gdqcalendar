package main

import (
	"log"

	parserapp "github.com/jcalmat/gdqcalendar/pkg/parser/app"
)

func main() {
	parserApp := parserapp.App{
		ScheduleURL:    "https://gamesdonequick.com/schedule/",
		ScheduleApiURL: "https://gamesdonequick.com/tracker/api/v2/events/{version}/runs?",
	}

	cal, err := parserApp.Parse()
	if err != nil {
		log.Fatalf("failed to parse gdq schedule: %s", err.Error())
	}

	err = cal.ToICS()
	if err != nil {
		log.Fatalf("failed to create ICS file: %s", err.Error())
	}
}
