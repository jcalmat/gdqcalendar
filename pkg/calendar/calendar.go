package calendar

import (
	"fmt"
	"os"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	guuid "github.com/google/uuid"
)

type Calendar struct {
	Name  string
	Games Games
}

type Games []Game

type Game struct {
	Name          string
	Category      string
	StartDate     time.Time
	SetupDuration time.Duration
	Duration      time.Duration
	Runners       []string
	Host          string
}

func (cal Calendar) ToICS() error {
	icsCal := ics.NewCalendar()
	icsCal.SetMethod(ics.MethodRequest)
	icsCal.SetName(cal.Name)
	for i, g := range cal.Games {
		event := icsCal.AddEvent(guuid.New().String())
		event.SetSummary(g.Name)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(g.StartDate)
		if i < len(cal.Games)-1 {
			event.SetEndAt(cal.Games[i+1].StartDate)
		} else {
			event.SetEndAt(g.StartDate.Add(g.Duration))
		}
		desc := fmt.Sprintf("Runners: %s\nRun Time: %s\nCategory: %s\nHost: %s\nSetup Time: %s", strings.Join(g.Runners, ", "), fmtDuration(g.Duration), g.Category, g.Host, fmtDuration(g.SetupDuration))
		event.SetDescription(desc)
	}

	// store it somewhere or return it idk yet
	calstr := icsCal.Serialize()

	f, err := os.Create("gdq.ics")
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(calstr)
	if err != nil {
		return err
	}

	err = f.Sync()
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
