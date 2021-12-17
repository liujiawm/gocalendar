package gocalendar

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDefaultCalendar(t *testing.T) {
	c := DefaultCalendar()
	t.Log(c.config.Grid)
}

func TestNewCalendar(t *testing.T) {
	cfg := CalendarConfig{
		Grid:1,
		TimeZoneName:"Europe/Berlin",
	}

	c := NewCalendar(cfg)

	ti := c.GetRawTime()

	t.Log(ti.Zone())
	t.Log(ti.String())
}

// 一整年节气
func TestCalendar_SolarTerms(t *testing.T) {
	c := NewCalendar(CalendarConfig{TimeZoneName:"Asia/Shanghai"})

	sts:= c.SolarTerms(2021)
	for _,v := range sts{
		st := fmt.Sprintf(" %s 定%s:%s", v.Name, v.Name, v.Time.Format(time.RFC3339))
		t.Log(st)
	}
}

func TestCalendar_GenerateWithDate(t *testing.T) {
	beforeTime := time.Now()
	defer func() {
		str := time.Since(beforeTime)
		t.Logf("本次执行用时：%s\n",str)
	}()

	c := NewCalendar(CalendarConfig{
		Grid:GridWeek,
		FirstWeek:0,
		SolarTerms:true,
		Lunar:true,
		HeavenlyEarthly:true,
		NightZiHour:true,
		StarSign:true,
	})
	result := c.GenerateWithDate(2021,12,22)
	for _,item := range result {
		fmt.Println(item)
	}
	t.Log("----------------------------")
	items := c.NextMonth()
	for _,item := range items {
		fmt.Println(item)
	}
}

func TestCalendar_Generate(t *testing.T) {
	beforeTime := time.Now()
	defer func() {
		str := time.Since(beforeTime)
		t.Logf("本次执行用时：%s\n",str)
	}()

	c := NewCalendar(CalendarConfig{
		Grid:GridMonth,
		FirstWeek:0,
		SolarTerms:true,
		Lunar:true,
		HeavenlyEarthly:true,
		NightZiHour:true,
		StarSign:true,
	})

	c.SetRawTime(2021,2,1)
	result := c.Generate()


	for _,item := range result {
		// t.Log(item)
		fmt.Println(item)
	}
}

// 一周日历表
func TestCalendar_weekCalendar(t *testing.T) {
	beforeTime := time.Now()
	defer func() {
		str := time.Since(beforeTime)
		t.Logf("本次执行用时：%s\n",str)
	}()

	c := NewCalendar(CalendarConfig{Grid:GridWeek,FirstWeek:0,SolarTerms:true,Lunar:true,HeavenlyEarthly:true,NightZiHour:true})
	c.weekCalendar()

	for _,item := range c.Items {
		t.Log(item)
	}
}

// 一月日历表
func TestCalendar_monthCalendar(t *testing.T) {
	beforeTime := time.Now()
	defer func() {
		str := time.Since(beforeTime)
		t.Logf("本次执行用时：%s\n",str)
	}()

	c := NewCalendar(CalendarConfig{Grid:GridMonth,FirstWeek:0,TimeZoneName:"Asia/Shanghai",SolarTerms:true,Lunar:true,HeavenlyEarthly:true,NightZiHour:true,StarSign:true})

	c.SetRawTime(2021,12,1)
	items := c.monthCalendar()


	itemsJson,_ := json.Marshal(items)
	t.Log(string(itemsJson))

	for _,item := range items {
		t.Log(item)
	}


	t.Logf("st len=%d jSS len=%d qSS len=%d tNM len=%d lMC len=%d lMD len=%d lFD len=%d gFD len=%d",len(c.tempData.st.data),len(c.tempData.jSS.data),len(c.tempData.qSS.data),len(c.tempData.tNM.data),len(c.tempData.lMC.data),len(c.tempData.lMD.data),len(c.tempData.lFD.data),len(c.tempData.gFD.data))
}

// 星座
func TestStarSign(t *testing.T) {
	i,ss,_ := StarSign(5,6)

	if i == 3 {
		t.Log("passed")
	}else{
		t.Error(i,ss)
	}
}

// Benchmark
func BenchmarkCalendar_monthCalendar(b *testing.B) {
	c := DefaultCalendar()
	c.monthCalendar()
}

// Benchmark
func BenchmarkCalendar_SolarTerms(b *testing.B) {
	c := NewCalendar(CalendarConfig{TimeZoneName:"Asia/Shanghai"})
	c.SolarTerms(2021)
}


