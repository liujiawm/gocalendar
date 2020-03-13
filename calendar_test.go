package gocalendar

import (
	"testing"
)

var c = NewCtime()
var d = CtimeNewDate(c)

// 24节气(从立春开始)
func TestDate_Jieqi(t *testing.T) {

	t.Log(NewCtime().Year(),"年从立春开始的节气时间：")
	jq := d.Jieqi()
	for k,v := range(jq){
		t.Logf("%d:[%s] %d-%d-%d %d:%d:%d\n",k,v.Name,v.Date.Year,v.Date.Month,v.Date.Day,v.Date.Hour,v.Date.Min,v.Date.Sec)
	}
}

// 公历的天干地支
func TestDate_GanZhi(t *testing.T) {
	/*d = NewDate(&Date{
		Year:2020,
		Month:2,
		Day:4,
		Hour:20,
		Min:35,
		Sec:42,
		Loc:NewCtime().Location(),
	})*/
	Zwz = true
	dc,_ := d.GanZhi()
	t.Logf("%s%s年 %s%s月 %s%s日 %s%s时\n",
		TianGanArray[dc.TianGanDiZhi.Ytg], DiZhiArray[dc.TianGanDiZhi.Ydz],
		TianGanArray[dc.TianGanDiZhi.Mtg],DiZhiArray[dc.TianGanDiZhi.Mdz],
		TianGanArray[dc.TianGanDiZhi.Dtg],DiZhiArray[dc.TianGanDiZhi.Ddz],
		TianGanArray[dc.TianGanDiZhi.Htg],DiZhiArray[dc.TianGanDiZhi.Hdz])
	// t.Logf("%#v\n",dc.TianGanDiZhi)
}

// 公历转农历，农历转公历
func TestDate_Solar2Lunar(t *testing.T) {
	d := NewDate(&Date{
		Year:2020,
		Month:6,
		Day:15,
		Hour:5,
		Min:35,
		Sec:42,
		Loc:NewCtime().Location(),
	})
	l := d.Solar2Lunar()
	leapMonthStr := ""
	if l.LeapMonth == 1 {
		leapMonthStr = "(闰)"
	}
	t.Logf("公历%d年%d月%d日转成农历是:%d年%s%s月%s日",d.Year,d.Month,d.Day,l.Year,leapMonthStr,MonthChinese(l.Month),DayChinese(l.Day))

	sd := l.Lunar2Solar()
	t.Logf("农历:%d年%s%s月%s日转成公历是:%d年%d月%d日",l.Year,leapMonthStr,MonthChinese(l.Month),DayChinese(l.Day),sd.Year,sd.Month,sd.Day)
}
// 星座
func TestDate_XingZuo(t *testing.T) {
	t.Logf("公历%d年%d月%d日是%s座",d.Year,d.Month,d.Day,d.XingZuo())
}
