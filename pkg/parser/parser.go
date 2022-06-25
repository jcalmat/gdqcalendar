package parser

import (
	"github.com/jcalmat/gdqcalendar/pkg/calendar"
)

type App interface {
	Parse() (calendar.C, error)
}
