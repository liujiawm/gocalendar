// 年数在1000-3000

package gocalendar

import (
	"testing"
)


func TestPureJQsinceSpring(t *testing.T){
	j := pureJQsinceSpring(2020)
	for _,v := range j {
		t.Log(v)
	}

}

var c = NewCtime()
var d = CtimeNewDate(c)

// 24节气(从立春开始)
func TestDate_Jieqi(t *testing.T) {

	t.Log(NewCtime().Year(),"年从立春开始的节气时间：")
	jq := d.Jieqi()
	for k,v := range(jq){
		t.Logf("%s:[%s] %d-%d-%d %d:%d:%d\n",k,v.Name,v.Date.Year,v.Date.Month,v.Date.Day,v.Date.Hour,v.Date.Min,v.Date.Sec)
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
	d.GanZhi()
	t.Logf("公历 %d年%d月%d日%d时 转成以立春开始的干支是: %s%s年 %s%s月 %s%s日 %s%s时\n",d.Year,d.Month,d.Day,d.Hour,
		TianGanArray[d.TianGanDiZhi.Ytg], DiZhiArray[d.TianGanDiZhi.Ydz],
		TianGanArray[d.TianGanDiZhi.Mtg],DiZhiArray[d.TianGanDiZhi.Mdz],
		TianGanArray[d.TianGanDiZhi.Dtg],DiZhiArray[d.TianGanDiZhi.Ddz],
		TianGanArray[d.TianGanDiZhi.Htg],DiZhiArray[d.TianGanDiZhi.Hdz])
	// t.Logf("%#v\n",dc.TianGanDiZhi)
}

// 公历转农历，农历转公历
func TestDate_Solar2Lunar(t *testing.T) {
	d := NewDate(&Date{
		Year:2020,
		Month:5,
		Day:28,
		Hour:0,
		Min:0,
		Sec:0,
		Loc:NewCtime().Location(),
	})
	l := d.Solar2Lunar()
	leapMonthStr := ""
	if l.LeapMonth == 1 {
		leapMonthStr = "(闰)"
	}
	t.Logf("公历%d年%d月%d日转成农历是:%d(%s%s)[%s]年%s%s月%s日",d.Year,d.Month,d.Day,l.Year,l.YearGanZi.Gan,l.YearGanZi.Zhi,l.YearGanZi.Animals, leapMonthStr,MonthChinese(l.Month),DayChinese(l.Day))

	sd := l.Lunar2Solar()
	t.Logf("农历:%d年%s%s月%s日转成公历是:%d年%d月%d日",l.Year,leapMonthStr,MonthChinese(l.Month),DayChinese(l.Day),sd.Year,sd.Month,sd.Day)
}
// 星座
func TestDate_XingZuo(t *testing.T) {
	t.Logf("公历%d年%d月%d日是%s座",d.Year,d.Month,d.Day,d.XingZuo())
}
