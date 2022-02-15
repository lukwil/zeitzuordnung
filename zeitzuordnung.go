package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/de"
)

const (
	startTime      = "09:00"
	endTime        = "18:00"
	reducedEndTime = "17:00"
	breakTime      = "01:00"
)

var records = [][]string{
	{"Datum", "Start", "Ende", "Pause"},
}

var c = cal.NewBusinessCalendar()

func main() {
	c.AddHoliday(de.HolidaysBY...)
	friedensfest := &cal.Holiday{
		Name:  "Friedensfest",
		Type:  cal.ObservancePublic,
		Month: time.August,
		Day:   8,
		Func:  cal.CalcDayOfMonth,
	}
	c.AddHoliday(friedensfest)
	c.AddHoliday(de.MariaHimmelfahrt)

	timezone, _ := time.LoadLocation("Europe/Berlin")
	today := time.Now().In(timezone)

	year, week := today.ISOWeek()
	fileName := fmt.Sprintf("zeitzuordnung-%v-KW%v.csv", year, week)
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	generateDataForWeek(today)

	w := csv.NewWriter(f)
	w.Comma = ';'
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	fmt.Printf("%v-KW%v", year, week)
}

func generateDataForWeek(date time.Time) {
	// starting from current day until friday
	for date.Weekday() != time.Saturday {
		if holiday, _, _ := c.IsHoliday(date); !holiday {
			hasReducedEndTime := false
			if date.Weekday() == time.Friday {
				// One hour less for friday --> 39h/week
				hasReducedEndTime = true
			}
			appendRecord(date, hasReducedEndTime)
		} else {
			if date.Weekday() == time.Friday {
				// When friday is a holiday, time can be reduced to seven hours on another day
				reduceTimeOnDayBefore()
			}
		}
		date = date.AddDate(0, 0, 1)
	}
}

func appendRecord(date time.Time, hasReducedEndTime bool) {
	r := []string{
		date.Format("2006-01-02"),
		startTime,
		endTime,
		breakTime,
	}

	if hasReducedEndTime {
		r[2] = reducedEndTime
	}
	records = append(records, r)
}

func reduceTimeOnDayBefore() {
	lastElement := records[len(records)-1]
	lastElement[2] = reducedEndTime
}
