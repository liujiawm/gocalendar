package gocalendar

import (
	"testing"
)

// func TestNewConfig(t *testing.T) {
// 	m := map[string]interface{}{
// 		"Grid":-1,
// 		"FirstWeek":-771,
// 		"TimeZoneName":"Asia/Chongqing",
// 		"SolarTerms":false,
// 		"Lunar":false,
// 		"HeavenlyEarthly":false,
// 		"NightZiHour":false,
// 	}
// 	cfg := NewConfig(m)
// 	t.Log(cfg)
// }

func TestCalendar_GetConfig(t *testing.T) {
	// c := DefaultCalendar()
	c := NewCalendar(CalendarConfig{Grid:-2,FirstWeek:789,TimeZoneName:"Asia/Chongqing",Lunar:true})
	cfg := c.GetConfig()
	t.Log(cfg)
}
