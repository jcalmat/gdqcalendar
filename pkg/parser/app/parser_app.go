package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jcalmat/gdqcalendar/pkg/calendar"
	"github.com/jcalmat/gdqcalendar/pkg/parser"
)

type App struct {
	ScheduleURL    string
	ScheduleApiURL string
}

func (a App) Parse() (calendar.C, error) {
	cal := calendar.C{
		Games: make(calendar.Games, 0),
	}

	// retrive redirect url from the schedule page
	apiURL := a.getApiURL()

	resp, err := http.Get(apiURL)
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

func (a App) getApiURL() string {
	var scheduleVersion string

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			scheduleVersion = strings.TrimPrefix(req.URL.String(), a.ScheduleURL)
			return nil
		},
	}

	resp, err := client.Get(a.ScheduleURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	return strings.Replace(a.ScheduleApiURL, "{version}", scheduleVersion, 1)
}
