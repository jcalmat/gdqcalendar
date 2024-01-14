package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jcalmat/gdqcalendar/pkg/calendar"
	"github.com/jcalmat/gdqcalendar/pkg/parser"
)

type App struct {
	ScheduleURL string
}

func (a App) Parse() (calendar.C, error) {
	cal := calendar.C{
		Games: make(calendar.Games, 0),
	}

	resp, err := http.Get(a.ScheduleURL)
	if err != nil {
		return cal, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return cal, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	var parsed parser.Calendar
	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return cal, err
	}

	for _, r := range parsed.Results {
		var game calendar.Game
		game.Name = r.Name
		game.Category = r.Category
		game.Runners = make([]string, 0)
		for _, runner := range r.Runners {
			game.Runners = append(game.Runners, runner.Name)
		}

		if len(r.Hosts) > 0 {
			pronouns := ""
			if r.Hosts[0].Pronouns != "" {
				pronouns = fmt.Sprintf(" (%s)", r.Hosts[0].Pronouns)
			}
			game.Host = fmt.Sprintf("%s%s", r.Hosts[0].Name, pronouns)
		}

		game.StartDate = r.StartTime
		game.EndDate = r.EndTime
		game.Duration = r.RunTime
		game.SetupDuration = r.SetupTime
		cal.Games = append(cal.Games, game)
	}

	return cal, nil
}
