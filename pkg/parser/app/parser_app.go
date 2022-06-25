package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jcalmat/gdqcalendar/pkg/calendar"
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
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return cal, err
	}

	// Find the title (should only be one)
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		cal.Name = s.First().Text()
	})

	doc.Find("tbody").First().Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			var game calendar.Game
			if i%2 == 0 {
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					switch i {
					case 0:
						game.StartDate, _ = time.Parse(time.RFC3339, s.Text())
						game.StartDate = game.StartDate.Local()
					case 1:
						game.Name = s.Text()
					case 2:
						game.Runners = strings.Split(s.Text(), ", ")
					case 3:
						game.SetupDuration, err = parseDuration(strings.TrimSpace(s.Text()))
						if err != nil {
							log.Println(err.Error())
						}
					}
				})
				cal.Games = append(cal.Games, game)
			} else {
				game = cal.Games[len(cal.Games)-1]
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					switch i {
					case 0:
						game.Duration, err = parseDuration(strings.TrimSpace(s.Text()))
						if err != nil {
							log.Println(err.Error())
						}
					case 1:
						game.Category = s.Text()
					case 2:
						game.Host = strings.TrimSpace(s.Text())
					}
				})
				cal.Games[len(cal.Games)-1] = game
			}
		})
	})

	return cal, nil
}

func parseDuration(s string) (time.Duration, error) {
	durations := strings.Split(s, ":")
	if len(durations) != 3 {
		return 0, errors.New("invalid duration")
	}

	hours, err := strconv.Atoi(durations[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(durations[1])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(durations[2])
	if err != nil {
		return 0, err
	}

	// nanoseconds conversion
	return time.Duration(hours*3600000000000 + minutes*60000000000 + seconds*1000000000), nil
}
