package main

import (
	"github.com/jcalmat/gdqcalendar/pkg/parser"
	"google.golang.org/api/calendar/v3"
)

type Gcal struct {
	Service    *calendar.Service
	CalendarID string
	GDQParser  parser.App
}
