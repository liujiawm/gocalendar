package gocalendar

import (
	"testing"
)

func TestGetMonthCalendar(t *testing.T) {
	c := GetMonthCalendar(2020,1)

	t.Log(c.Year,"年",c.Month,"月")
	t.Log("今天：",c.ToDay)
	for _,v := range c.Items {
		leapMonthStr := ""
		if v.LunarDate.LeapMonth == 1 {
			leapMonthStr = "(闰)"
		}
		t.Logf("公历年:%d年%d月%d日周%d 农历年:%d(%s%s)[%s]年%s%s月%s日 节气:%s",v.Date.Year,v.Date.Month,v.Date.Day,v.Date.Week,
			v.LunarDate.Year, v.LunarDate.YearGanZi.Gan, v.LunarDate.YearGanZi.Zhi, v.LunarDate.YearGanZi.Animals,
			leapMonthStr,MonthChinese(v.LunarDate.Month),DayChinese(v.LunarDate.Day),v.Date.JQ.Name)
	}

	//t.Log(nd.Solar2Lunar())
}
