package gocalendar

import (
	"fmt"
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
		jqstr := ""
		if v.JQ != nil {
			jqstr = fmt.Sprintf(" 节气:%s (定%s：%d:%d:%d) \n",v.JQ.Name,v.JQ.Name,v.JQ.Date.Hour,v.JQ.Date.Min,v.JQ.Date.Sec)
		}

		t.Logf("公历年:%d年%d月%d日周%d 农历年:%d(%s%s)[%s]年%s%s月%s日%s",v.Date.Year,v.Date.Month,v.Date.Day,v.Date.Week,
			v.LunarDate.Year, v.LunarDate.YearGanZi.Gan, v.LunarDate.YearGanZi.Zhi, v.LunarDate.YearGanZi.Animals,
			leapMonthStr,MonthChinese(v.LunarDate.Month),DayChinese(v.LunarDate.Day),jqstr)
	}

	//t.Log(nd.Solar2Lunar())
}
