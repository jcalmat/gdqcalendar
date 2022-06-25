package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	CalendarID string `json:"calendarId"`
}

func parseConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
