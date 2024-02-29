package types

import (
	"fmt"
	"time"
)

type Time time.Time

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t Time) IsSameDate(tm time.Time) bool {
	compareTime := t.Time()

	if tm.Day() != compareTime.Day() {
		return false
	}

	if tm.Month() != compareTime.Month() {
		return false
	}

	if tm.Year() != compareTime.Year() {
		return false
	}

	return true
}

func (t Time) BeginningOfMonth() time.Time {
	y, m, _ := t.Time().Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

func (t Time) EndOfMonth() time.Time {
	return t.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond).Round(time.Millisecond)
}

func (t Time) ToCCYYMMDD() string {
	location, _ := time.LoadLocation("America/New_York")
	localTime := t.Time().In(location)

	return fmt.Sprintf("%04d%02d%02d", localTime.Year(), localTime.Month(), localTime.Day())
}

func (t Time) ToCCYYMMDD_HHMMSS() string {
	location, _ := time.LoadLocation("America/New_York")
	localTime := t.Time().In(location)

	return fmt.Sprintf("%04d%02d%02d %02d%02d%02d",
		localTime.Year(), localTime.Month(), localTime.Day(),
		localTime.Hour(), localTime.Minute(), localTime.Second())
}

func (t *Time) ConvertCCYYMMDDFormatToTime(timeString string) Time {
	layout := "2006-01-2 15:04:05"
	location, _ := time.LoadLocation("America/New_York")
	newTimeStruct, _ := time.ParseInLocation(layout, timeString, location)
	*t = Time(newTimeStruct.UTC())

	return *t
}

func (t *Time) ConvertCCYYMMDDFormatToTickTime(timeString string) Time {
	layout := "2006-01-2 15:04:05.999999"
	location, _ := time.LoadLocation("America/New_York")
	newTimeStruct, _ := time.ParseInLocation(layout, timeString, location)
	*t = Time(newTimeStruct.UTC())

	return *t
}

func (t *Time) ConvertCCYYMMDDDailyFormatToTime(timeString string) Time {
	layout := "2006-01-2"
	location, _ := time.LoadLocation("America/New_York")
	newTimeStruct, _ := time.ParseInLocation(layout, timeString, location)
	*t = Time(newTimeStruct.UTC())

	return *t
}
