package main

import (
	"context"
	"embed"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type CalendarDate struct {
	Year  int
	Month int
	Day   int
}

func NewCalendarDate(date time.Time) *CalendarDate {
	return &CalendarDate{
		Year:  date.Year(),
		Month: int(date.Month()),
		Day:   date.Day(),
	}
}

type CalendarCell struct {
	Date      *CalendarDate
	ClassName string
	Label     string
}

func NewCalendarCell(date time.Time, className string, label string) *CalendarCell {
	return &CalendarCell{
		Date:      NewCalendarDate(date),
		ClassName: className,
		Label:     label,
	}
}

type Calendar struct {
	Current *CalendarCell
	Days    []*CalendarCell
}

func NewCalendar(current time.Time) *Calendar {
	firstDay := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())
	topLeftDay := firstDay.AddDate(0, 0, int(firstDay.Weekday())*(-1))
	lastDay := firstDay.AddDate(0, 1, -1)

	sum := int(firstDay.Weekday()) + lastDay.Day()
	rowCount := sum/7 + 1
	if sum%7 == 0 {
		rowCount--
	}

	calendar := &Calendar{}

	for i := 0; i < rowCount*7; i++ {
		d := topLeftDay.AddDate(0, 0, i)

		className := "calendar-cell-weekday"
		if d.Weekday() == time.Sunday {
			className = "calendar-cell-sunday"
		} else if d.Weekday() == time.Saturday {
			className = "calendar-cell-saturday"
		}

		label := ""
		if current.Month() == d.Month() {
			label = strconv.Itoa(d.Day())
		}

		cell := NewCalendarCell(d, className, label)
		calendar.Days = append(calendar.Days, cell)

		if d.Year() == current.Year() &&
			d.Month() == current.Month() &&
			d.Day() == current.Day() {
			calendar.Current = cell
			cell.ClassName = cell.ClassName + " calendar-cell-today"
		}
	}

	return calendar
}

type App struct {
	ctx     context.Context
	current time.Time
}

func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.current = time.Now()
}

func (a *App) GetCalendar() *Calendar {
	return NewCalendar(a.current)
}

func (a *App) SetCalendar(year, month, day int) *Calendar {
	a.current = time.Date(year, time.Month(month), day, 0, 0, 0, 0, a.current.Location())
	return a.GetCalendar()
}

func (a *App) OffsetCalendar(year, month, day int, isWrapDay bool) *Calendar {
	if isWrapDay {
		a.current = a.current.AddDate(year, month, day)
	} else {
		new := time.Date(a.current.Year(), a.current.Month(), 1, 0, 0, 0, 0, a.current.Location())
		new = new.AddDate(year, month, 0)
		newLastDay := new.AddDate(0, 1, -1)

		d := a.current.Day() + day
		d = min(d, newLastDay.Day())
		d = max(d, 1)

		new = new.AddDate(0, 0, d-1)
		a.current = new
	}

	return a.GetCalendar()
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "カレンダー",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.Startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
