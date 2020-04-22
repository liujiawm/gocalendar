package gocalendar

import (
	"strings"
	"testing"
	"time"
)

func TestNewCtime(t *testing.T) {
	ct := NewCtime()
	t.Log(ct.Format(time.RFC3339))

	ct = UnixNewCtime(ct.Unix(),0)
	t.Log(ct.Format(time.RFC3339))

	ct = DateNewCtime(&Date{
		Year:2020,
		Month:2,
		Day:90,
	})
	t.Log(ct.Format(time.RFC3339))
}

func TestCtime_MonthDay(t *testing.T) {
	t1, _ := time.Parse("2006-01-02 15:04:05", "1583-10-12 00:00:00")
	t.Log(t1.Format(time.RFC3339))
	ct := Ctime(t1)
	t.Log(ct.MonthDay())
	a := "2020-03-25"
	as := strings.SplitN(a,"-",4)
	t.Log(as[0],as[1],as[2])
}