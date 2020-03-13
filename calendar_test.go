package gocalendar

import (
	"fmt"
	"testing"
)

var c = NewCtime()
var d = CtimeNewDate(c)

func TestDate_GanZhi(t *testing.T) {
	d = NewDate(&Date{
		Year:5020,
		Month:2,
		Day:4,
		Hour:20,
		Min:35,
		Sec:42,
		Loc:NewCtime().Location(),
	})
	fmt.Println("=======================================================================================================")
	Zwz = true
	dc,_ := d.GanZhi()
	fmt.Printf("%#v\n",dc.TianGanDiZhi)
}
func TestDate_Jieqi(t *testing.T) {

	d := NewDate(&Date{
		Year:2020,
		Month:6,
		Day:15,
		Hour:5,
		Min:35,
		Sec:42,
		Loc:NewCtime().Location(),
	})

	fmt.Println("=======================================================================================================")
	fmt.Println("2020年从立春开始的节气时间：")
	jq := d.Jieqi()
	for k,v := range(jq){
		fmt.Printf("%d:[%s] %d-%d-%d %d:%d:%d\n",k,v.Name,v.Date.Year,v.Date.Month,v.Date.Day,v.Date.Hour,v.Date.Min,v.Date.Sec)
	}


}
