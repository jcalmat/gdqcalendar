package calendar

import (
	"fmt"
	"os"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type C struct {
	Name  string
	Games Games
}

type Games []Game

type Game struct {
	Name          string
	Category      string
	StartDate     time.Time
	EndDate       time.Time
	SetupDuration string
	Duration      string
	Runners       []string
	Host          string
}

func (cal C) ToICS() error {
	icsCal := ics.NewCalendar()
	icsCal.SetMethod(ics.MethodRequest)
	icsCal.SetName(cal.Name)
	for i, g := range cal.Games {
		event := icsCal.AddEvent(fmt.Sprintf("%s_game_%02d", strings.ToLower(strings.ReplaceAll(cal.Name, " ", "_")), i))
		event.SetSummary(g.Name)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(g.StartDate)
		if i < len(cal.Games)-1 {
			event.SetEndAt(cal.Games[i+1].StartDate)
		} else {
			event.SetEndAt(cal.Games[i].EndDate)
		}
		desc := fmt.Sprintf("Runners: %s\nRun Time: %s\nCategory: %s\nHost: %s\nSetup Time: %s", strings.Join(g.Runners, ", "), g.Duration, g.Category, g.Host, g.SetupDuration)
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
