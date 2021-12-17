package gocalendar

import (
	// "math"
	// "reflect"
	"time"
)

// 取日历方式
const (
	GridDay int = iota
	GridWeek
	GridMonth
)

// type CalendarConfig struct 配置
type CalendarConfig struct {
	Grid            int    // 取日历方式,GridDay按天取日历,GridWeek按周取日历,GridMonth按月取日历
	FirstWeek       int    // 日历显示时第一列显示周几，(日历表第一列是周几,0周日,依次最大值6)
	TimeZoneName    string // 时区名称,需php支持的时区名称
	SolarTerms      bool   // 读取节气 bool
	Lunar           bool   // 读取农历 bool
	HeavenlyEarthly bool   // 读取干支 bool
	NightZiHour     bool   // 区分早晚子时，true则 23:00-24:00 00:00-01:00为子时，否则00:00-02:00为子时
	StarSign        bool   // 读取星座
}

// defaultConfig 新的默认配置
func defaultConfig() CalendarConfig {
	return CalendarConfig{
		Grid:            GridMonth,
		FirstWeek:       0,
		TimeZoneName:    time.Local.String(),
		SolarTerms:      true,
		Lunar:           true,
		HeavenlyEarthly: true,
		NightZiHour:     true,
		StarSign:        true,
	}
}

// (*CalendarConfig) clone
func (cfg *CalendarConfig) clone() *CalendarConfig {
	return &CalendarConfig{
		Grid:            cfg.Grid,
		FirstWeek:       cfg.FirstWeek,
		TimeZoneName:    cfg.TimeZoneName,
		SolarTerms:      cfg.SolarTerms,
		Lunar:           cfg.Lunar,
		HeavenlyEarthly: cfg.HeavenlyEarthly,
		NightZiHour:     cfg.NightZiHour,
		StarSign:        cfg.StarSign,
	}
}

// (*Calendar) GetConfig 读取配置
// 返回的是一个clone
func (c *Calendar) GetConfig() CalendarConfig {
	cfg := c.config.clone()
	return *cfg
}

// NewConfig 用一个map[string]interface{}新建配置
// func NewConfig(cfgMap map[string]interface{}) CalendarConfig {
// 	cfg := defaultConfig()
//
// 	if grid, ok := cfgMap["Grid"]; ok {
// 		if reflect.TypeOf(grid).Kind() == reflect.Int {
//
// 			cfg.Grid = int(math.Mod(math.Abs(float64(reflect.ValueOf(grid).Int())),3))
// 		}
// 	}
//
// 	if firstWeek, ok := cfgMap["FirstWeek"]; ok {
// 		if reflect.TypeOf(firstWeek).Kind() == reflect.Int {
// 			cfg.FirstWeek = int(math.Mod(math.Abs(float64(reflect.ValueOf(firstWeek).Int())),7))
// 		}
// 	}
//
// 	if timeZoneName, ok := cfgMap["TimeZoneName"]; ok {
// 		if reflect.TypeOf(timeZoneName).Kind() == reflect.String {
// 			tzn := reflect.ValueOf(timeZoneName).String()
// 			if tzn != "" {
// 				cfg.TimeZoneName = reflect.ValueOf(timeZoneName).String()
// 			}
// 		}
// 	}
//
// 	if solarTerms, ok := cfgMap["SolarTerms"]; ok {
// 		if reflect.TypeOf(solarTerms).Kind() == reflect.Bool {
// 			cfg.SolarTerms = reflect.ValueOf(solarTerms).Bool()
// 		}
// 	}
//
// 	if lunar, ok := cfgMap["Lunar"]; ok {
// 		if reflect.TypeOf(lunar).Kind() == reflect.Bool {
// 			cfg.Lunar = reflect.ValueOf(lunar).Bool()
// 		}
// 	}
//
// 	if heavenlyEarthly, ok := cfgMap["HeavenlyEarthly"]; ok {
// 		if reflect.TypeOf(heavenlyEarthly).Kind() == reflect.Bool {
// 			cfg.HeavenlyEarthly = reflect.ValueOf(heavenlyEarthly).Bool()
// 		}
// 	}
//
// 	if nightZiHour, ok := cfgMap["NightZiHour"]; ok {
// 		if reflect.TypeOf(nightZiHour).Kind() == reflect.Bool {
// 			cfg.NightZiHour = reflect.ValueOf(nightZiHour).Bool()
// 		}
// 	}
//
// 	if starSign, ok := cfgMap["StarSign"]; ok {
// 		if reflect.TypeOf(starSign).Kind() == reflect.Bool {
// 			cfg.StarSign = reflect.ValueOf(starSign).Bool()
// 		}
// 	}
//
// 	return cfg
// }


