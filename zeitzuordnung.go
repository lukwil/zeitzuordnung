package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

var records = [][]string{
	{"Datum", "Start", "Ende", "Pause"},
}

func main() {
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
	// going back from current day until monday
	startTime := "09:00"
	endTime := "18:00"
	breakTime := "01:00"
	for date.Weekday() != time.Monday {
		appendEntry(date, startTime, endTime, breakTime)

		date = date.AddDate(0, 0, -1)
	}

	// One hour less for monday --> 39h/week
	endTime = "17:00"
	appendEntry(date, startTime, endTime, breakTime)
}

func appendEntry(date time.Time, startTime, endTime, breakTime string) {
	records = append(records, []string{
		date.Format("02.01.2006"),
		startTime,
		endTime,
		breakTime,
	})
}
