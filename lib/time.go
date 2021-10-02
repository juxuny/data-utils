package lib

import "time"

type _Time struct{}

var Time = _Time{}

func (_Time) NowPointer() *time.Time {
	t := time.Now()
	return &t
}

func (_Time) Parse(layout string, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

const DayLayout = "2006-01-02"
const TimeLayout = "15:04:05"
const DateTimeLayout = DayLayout + " " + TimeLayout
const TimeInMillionLayout = TimeLayout + ".000"
