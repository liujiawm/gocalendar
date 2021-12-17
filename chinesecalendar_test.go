package gocalendar

import (
	"fmt"
	"testing"
	"time"
)

// 干支
func TestCalendar_ChineseSexagenaryCycle(t *testing.T) {
	c := NewCalendar(CalendarConfig{NightZiHour: false})
	rt := time.Date(2021, 5, 6, 23, 50, 0, 0, time.Local)
	gz := c.ChineseSexagenaryCycle(rt)
	if gz.Year.HSI == 7 && gz.Year.EBI == 1 && gz.Month.HSI == 9 && gz.Month.EBI == 5 && gz.Day.HSI == 1 && gz.Day.EBI == 3 {
		t.Log("passed")
	} else {
		t.Error(gz.Year, gz.Month, gz.Day, gz.Hour)
	}

	c = NewCalendar(CalendarConfig{NightZiHour: true})
	rt = time.Date(2021, 5, 6, 23, 50, 0, 0, time.Local)
	gz = c.ChineseSexagenaryCycle(rt)
	if gz.Day.HSI == 0 && gz.Day.EBI == 2 {
		t.Log("passed")
	} else {
		t.Error(gz)
	}

}

// 公历转农历
func TestCalendar_GregorianToLunar(t *testing.T) {
	ld := DefaultCalendar().GregorianToLunar(1000, 6, 5)
	t.Log(ld)
}

// 农历转公历
func TestCalendar_LunarToGregorian(t *testing.T) {
	c := DefaultCalendar()
	gd, _ := c.LunarToGregorian(2020, 4, 14, false)
	fmt.Println("农历2020年四月十四转换成公历是:", gd.Format("2006-01-02"))

	gd, _ = c.LunarToGregorian(2020, 4, 14, true)
	fmt.Println("农历2020年闰四月十四转换成公历是:", gd.Format("2006-01-02"))

	t.Log(gd.Format(time.RFC3339))
}

func TestCalendar_lunarToGregorian(t *testing.T) {
	dc := DefaultCalendar()

	gd, _ := dc.lunarToGregorian(dc.GregorianToLunar(2020, 6, 5))
	t.Log(gd.Format(time.RFC3339))
}

// 农历某月天数
func TestCalendar_LunarMonthDay(t *testing.T) {
	dc := DefaultCalendar()
	d, e := dc.LunarMonthDays(2018, 12, false)

	if e != nil {
		t.Log(e.Error())
	} else {
		if d == 30 {
			t.Log("passed")
		} else {
			t.Error(d)
		}
	}
}

func TestLunarFestival(t *testing.T) {
	c := DefaultCalendar()
	lf := c.lunarFestival(2021, 12, 29, false)
	t.Log(lf)

	lf = c.lunarFestival(2021, 7, 7, false)
	t.Log(lf)
}

func TestCalendar_pureJieSinceSpring(t *testing.T) {
	dc := DefaultCalendar()

	year := 2021
	jss := dc.pureJieSinceSpring(year)

	for i, v := range jss {
		t.Logf("%d %.10f", i, v)
	}
}

func TestCalendar_qiSinceWinterSolstice(t *testing.T) {
	dc := DefaultCalendar()

	year := 2021
	qss := dc.qiSinceWinterSolstice(year)

	for i, v := range qss {
		t.Logf("%d %.10f", i, v)
	}
}
