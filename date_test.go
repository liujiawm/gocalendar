package gocalendar

import (
	"testing"
	"time"
)

func TestSolarDate_MonthDays(t *testing.T) {
	days := NewDate().monthDays()
	t.Log(days)
	days2 := YmdNewDate(2020,2,5,time.Local).monthDays()
	t.Log(days2)
	days3 := YmdNewDate(2019,2,5,time.Local).monthDays()
	t.Log(days3)
	days4 := YmdNewDate(2016,2,5,time.Local).monthDays()
	t.Log(days4)
	time.Now().Minute()
}

