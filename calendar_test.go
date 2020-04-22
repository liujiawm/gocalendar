package gocalendar

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// 日历:公历和农历,及json
func TestCalendar_Calendars(t *testing.T) {
	sc := Calendar{
		Loc:       time.Local, // 时区,默认time.Local
		FirstWeek: 0,          // 从周几天始,默认0周日开始, (日历表第列是周几,0周日,依次最大值6)
		Grid:      GridWeek,  // 取日历方式,默认取一个月, gocalendar.GridDay取一天,gocalendar.GridWeek 取一周,gocalendar.GridMonth取一个月
		Zwz:       false,      // 是否区分早晚子时(子时从23-01时),true则23:00-24:00算成上一天
		Getjq:     true,       // 是否取节气
	}
	cds := sc.Calendars(2020,4,22)
	for _,cd := range cds {
		var jqStr string
		if cd.SD.Jq != nil {
			jqStr = fmt.Sprintf(" 节气:%s (定%s：%d:%d:%d)",cd.SD.Jq.Name,cd.SD.Jq.Name,cd.SD.Jq.Date.Hour,cd.SD.Jq.Date.Min,cd.SD.Jq.Date.Sec)
		}

		t.Logf("公历:%d年%d月%d日 周%s 农历:%d(%s%s)[%s]年%s%s月%s  %s",cd.SD.Date.Year, cd.SD.Date.Month, cd.SD.Date.Day, WeekChinese(cd.SD.Date.Week), cd.LD.Date.Year,cd.LD.YearGanZi.Gan,cd.LD.YearGanZi.Zhi,cd.LD.YearGanZi.Animal, cd.LD.LeapStr, cd.LD.MonthStr, cd.LD.DayStr, jqStr)
	}
	cdjson,err := json.Marshal(cds)
	if err != nil {
		t.Error("JSON ERR:", err.Error())
	}
	t.Logf(string(cdjson))
}

// 日历:公历和农历
func BenchmarkCalendar_Calendars(b *testing.B) {
	sc := Calendar{
		Loc:       time.Local, // 时区,默认time.Local
		FirstWeek: 1,          // 从周几天始,默认0周日开始, (日历表第列是周几,0周日,依次最大值6)
		Grid:      GridMonth,  // 取日历方式,默认取一个月, gocalendar.GridDay取一天,gocalendar.GridWeek 取一周,gocalendar.GridMonth取一个月
		Zwz:       false,      // 是否区分早晚子时(子时从23-01时),true则23:00-24:00算成上一天
		Getjq:     true,       // 是否取节气
	}

	cds := sc.Calendars(2020,6,5)
	for _,cd := range cds {
		var jqStr string
		if cd.SD.Jq != nil {
			jqStr = fmt.Sprintf(" 节气:%s (定%s：%d:%d:%d)",cd.SD.Jq.Name,cd.SD.Jq.Name,cd.SD.Jq.Date.Hour,cd.SD.Jq.Date.Min,cd.SD.Jq.Date.Sec)
		}

		b.Logf("公历:%d年%d月%d日 周%s 农历:%d(%s%s)[%s]年%s%s月%s (%s%s年%s%s月%s%s日%s%s时) %s",
			cd.SD.Date.Year, cd.SD.Date.Month, cd.SD.Date.Day, WeekChinese(cd.SD.Date.Week),
			cd.LD.Date.Year,cd.LD.YearGanZi.Gan,cd.LD.YearGanZi.Zhi,cd.LD.YearGanZi.Animal, cd.LD.LeapStr, cd.LD.MonthStr, cd.LD.DayStr,
			cd.SD.GanZhi.YtgStr, cd.SD.GanZhi.YdzStr, cd.SD.GanZhi.MtgStr, cd.SD.GanZhi.MdzStr, cd.SD.GanZhi.DtgStr, cd.SD.GanZhi.DdzStr, cd.SD.GanZhi.HtgStr, cd.SD.GanZhi.HdzStr,
			jqStr)
	}
}

// 日历:公历
func TestCalendar_SolarCalendar(t *testing.T) {
	// 默认设置 DefaultCalendar()
	datas2 := DefaultCalendar().SolarCalendar(2020,12,31)
	if datas2 != nil {
		for _,d := range datas2 {
			var jqstr string
			if d.Jq != nil {
				jqstr = fmt.Sprintf(" 节气:%s (定%s：%d:%d:%d) \n",d.Jq.Name,d.Jq.Name,d.Jq.Date.Hour,d.Jq.Date.Min,d.Jq.Date.Sec)
			}
			t.Logf("%d年%d月%d日 %d时%d分%d秒 周%d (%s%s年%s%s月%s%s日%s%s时) %s",d.Year,d.Month,d.Day,d.Hour,d.Min,d.Sec,d.Week,
				d.GanZhi.YtgStr, d.GanZhi.YdzStr, d.GanZhi.MtgStr, d.GanZhi.MdzStr, d.GanZhi.DtgStr, d.GanZhi.DdzStr, d.GanZhi.HtgStr, d.GanZhi.HdzStr,
				jqstr)
		}
	}
}

// 节气
func TestCalendar_Jieqi(t *testing.T) {
	jq,_,err := DefaultCalendar().Jieqi(2020)
	if err != nil {
		t.Error(err.Error())
	}
	for _,d := range jq {
		t.Logf("%s - %d年%d月%d日 %d时%d分%d秒 周%d",d.Name,d.Date.Year,d.Date.Month,d.Date.Day,d.Date.Hour,d.Date.Min,d.Date.Sec,d.Date.Week)
	}
}

// 公历农历互换
func TestCalendar_Solar2Lunar(t *testing.T) {
	sc := Calendar{
		Getjq:     false,       // 是否取节气
	}

	d := YmdNewDate(2020,5,28,time.Local)
	ld,err := sc.Solar2Lunar(sc.DateToSolarDate(d))
	if err != nil {
		t.Error(err)
	}

	t.Logf("公历%d年%d月%d日转成农历是: %d(%s%s)[%s]年%s%s月%s",d.Year,d.Month,d.Day,ld.Year,ld.YearGanZi.Gan,ld.YearGanZi.Zhi,ld.YearGanZi.Animal, ld.LeapStr, ld.MonthStr, ld.DayStr)

	ssd,err := sc.Lunar2Solar(ld.Year,ld.Month,ld.Day,ld.LeapMonth)
	if err != nil {
		t.Error(err)
	}
	t.Logf("农历%d年%d月%d日转成公历是: %d年%d月%d日",ld.Year,ld.Month,ld.Day,ssd.Year,ssd.Month,ssd.Day)

}

// 星座
func TestZodiac(t *testing.T) {
	// 用日期的月和日取出对应的星座索引和名称
	i,zo := Zodiac(4,22)
	t.Log(i)
	t.Log(zo)
}