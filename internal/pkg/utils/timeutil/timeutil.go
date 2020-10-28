// Package timeutil contains all information that related to the time processing features
package timeutil

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
)

// Now return time now as a pointer
func Now() *time.Time {
	now := time.Now()
	return &now
}

// GetUnixMiliseconds return unix milicecond
func GetUnixMiliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// GetUnix is doing the formatting for time data
func GetUnix(time *time.Time) int64 {
	if time == nil {
		return 0
	}
	return time.Unix()
}

// OneYearLater return one year later from the given time
func OneYearLater(t time.Time) time.Time {
	return now.New(t.AddDate(1, 0, 0)).EndOfDay()
}

// OneMonthLater return one month later from the given time
func OneMonthLater(t time.Time) time.Time {
	return now.New(t.AddDate(0, 0, 30)).EndOfDay()
}

// OneWeekLater return one week later from the given time
func OneWeekLater(t time.Time) time.Time {
	return now.New(t.AddDate(0, 0, 7)).EndOfDay()
}

// ThaiDate return date in Thai format.
func ThaiDate(t time.Time) string {
	tm := []string{"ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}

	return fmt.Sprintf("%d %s %d", t.Day(), tm[int(t.Month())-1], t.Year()+543)
}

//Get last time of the month
func LastTimeOfTheMonth(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	year := t.Year()
	month := t.Month()

	nextMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := nextMonth.Day()
	expiredDay := time.Date(year, month, lastDayOfMonth, 16, 59, 59, 999999999, time.UTC)

	return &expiredDay
}
