package gocalendar

import (
	"fmt"
	"time"
)

// 日历显示时第一列显示周几，周日为0，默认周日
var FirstWeek = 0

type Calendars struct {
	Year      int
	Month     int
	ToDay     *Date
	Items     [42]Calendar
}
type Calendar struct {
	Date      *Date
	LunarDate *LunarDate
}

// 用年和月取一个月日历
func GetMonthCalendar(y,m int)*Calendars{
	y,m,_ = DateNewCtime(&Date{
		Year:y,
		Month:m,
		Day:1,
		Hour:0,
		Min:0,
		Sec:0,
		Nsec:0,
		Loc:time.Local,
	}).Date()

	if y < 1582 {
		y = 1582
	} else if y > 3000 {
		y = 3000
	}
	// 1582年11月之前的都不取了，太麻烦，有空再搞了
	if y==1582 && m <11 {
		m = 11
	}

	td := NewDate(&Date{
		Year:y,
		Month:m,
		Day:1,
		Hour:23,
		Min:59,
		Sec:59,
		Nsec:0,
		Loc:time.Local,
	})

	// 节气，map[string]*JQ 共27个值，从上一年的冬至开始
	jq := td.Jieqi()

	// 本月有多少天
	monthDays := td.solarDays()
	// 本月1日是周几
	tdfw := td.firstdayWeekday()
	// 本月最后一天是周几
	tdlw := td.lastdayWeekday()

	// 1日前空多少格
	prevFree := 7 + tdfw - FirstWeek // prevFree最大等于13
	if prevFree >= 7 {
		prevFree -= 7
	}
	// 最后日剩多少格
	nextFree := 6 - tdlw + FirstWeek // nextFree最大等于7

	// 后面补7天凑够42天
	if monthDays + prevFree + nextFree == 35 {
		nextFree += 7
	}

	var items [42]Calendar

	if prevFree > 0 {
		for i := 0; i < prevFree; i++ {
			nd := td.AddDate(0,0,-prevFree+i)
			// 农历
			nld := nd.Solar2Lunar()
			// 节气
			jqindex := fmt.Sprintf("%d-%d-%d",nd.Year,nd.Month,nd.Day)
			if jqv,ok := jq[jqindex];ok{
				nd.JQ = jqv
			}

			items[i] = Calendar{
				Date:nd,
				LunarDate:nld,
			}
		}
	}
	for i := 0; i < monthDays; i++ {
		nd := td.AddDate(0,0,i)
		// 农历
		nld := nd.Solar2Lunar()
		// 节气
		jqindex := fmt.Sprintf("%d-%d-%d",nd.Year,nd.Month,nd.Day)
		if jqv,ok := jq[jqindex];ok{
			nd.JQ = jqv
		}

		items[i+prevFree] = Calendar{
			Date:nd,
			LunarDate:nld,
		}
	}
	if nextFree > 0 {
		f := prevFree+monthDays
		fd := items[f-1].Date
		for i := 0; i < nextFree; i++ {
			nd := fd.AddDate(0,0,i+1)
			// 农历
			nld := nd.Solar2Lunar()
			// 节气
			jqindex := fmt.Sprintf("%d-%d-%d",nd.Year,nd.Month,nd.Day)
			if jqv,ok := jq[jqindex];ok{
				nd.JQ = jqv
			}

			items[i+f] = Calendar{
				Date:nd,
				LunarDate:nld,
			}
		}
	}

	today := CtimeNewDate(NewCtime())

	return &Calendars{
		Year:y,
		Month:m,
		ToDay:today,
		Items:items,
	}
}