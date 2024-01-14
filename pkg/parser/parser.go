package parser

import (
	"time"

	"github.com/jcalmat/gdqcalendar/pkg/calendar"
)

type App interface {
	Parse() (calendar.C, error)
}

type Calendar struct {
	Count   int       `json:"count"`
	Results []Results `json:"results"`
}
type Runners struct {
	Type     string `json:"type"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Pronouns string `json:"pronouns"`
}
type Results struct {
	Type         string        `json:"type"`
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	DisplayName  string        `json:"display_name"`
	Description  string        `json:"description"`
	Category     string        `json:"category"`
	Console      string        `json:"console"`
	Runners      []Runners     `json:"runners"`
	Hosts        []Host        `json:"hosts"`
	Commentators []interface{} `json:"commentators"`
	StartTime    time.Time     `json:"starttime"`
	EndTime      time.Time     `json:"endtime"`
	Order        int           `json:"order"`
	RunTime      string        `json:"run_time"`
	SetupTime    string        `json:"setup_time"`
}

type Host struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Pronouns string `json:"pronouns"`
	Type     string `json:"type"`
}
