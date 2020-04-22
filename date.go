package gocalendar

import (
	"time"
)

type Date struct {
	Year  int            `json:"year"`  // 年
	Month int            `json:"month"` // 月
	Day   int            `json:"day"`   // 日
	Hour  int            `json:"hour"`  // 时
	Min   int            `json:"min"`   // 分
	Sec   int            `json:"sec"`   // 秒
	Nsec  int            `json:"-"`     // 纳秒
	Week  int            `json:"week"`  // 星期几,0表示星期天
	Loc   *time.Location `json:"-"`     // 时区
}

var (
	// 当前时间函数
	timeFun = time.Now
)


// NewDate 用time.Local系统时区初始SolarDate
func NewDate()*Date{
	now := time.Now()
	return TimeToDate(now)
}

// UtcNewDate 用UTC时区初始化SolarDate
func UtcNewDate()*Date{
	return LocNewDate(time.UTC)
}

// ZoneNewDate 用timezone初始SolarDate
func ZoneNewDate(zone int)*Date{
	timeLoc := time.FixedZone(LocName(zone),zone)
	return LocNewDate(timeLoc)
}

// LocNewDate 用Location初始SolarDate
func LocNewDate(loc *time.Location)*Date{
	if loc == nil {
		loc = time.Local
	}
	now := timeFun().In(loc)
	return TimeToDate(now)
}

// YmdNewDate 指定年月日初始SolarDate
func YmdNewDate(y,m,d int,loc *time.Location)*Date{
	if loc == nil {
		loc = time.Local
	}
	t := time.Date(y,time.Month(m),d,0,0,0,0,loc)
	return TimeToDate(t)
}

// TimeToSolarDate Time转SolarDate
func TimeToDate(t time.Time)*Date{
	y,m,d := t.Date()
	h,mi,s := t.Clock()
	n := t.Nanosecond()
	w := t.Weekday()
	return &Date{
		Year:y,
		Month:int(m),
		Day:d,
		Hour:h,
		Min:mi,
		Sec:s,
		Nsec:n,
		Week:int(w),
		Loc:t.Location(),
	}
}

// clone
func (sd *Date)clone()*Date{
	return &Date{
		Year:sd.Year,
		Month:sd.Month,
		Day:sd.Day,
		Hour:sd.Hour,
		Min:sd.Min,
		Sec:sd.Sec,
		Nsec:sd.Nsec,
		Week:sd.Week,
		Loc:sd.Loc,
	}
}

// add 加 years年,months月,days日
func (sd *Date)add(years, months, days int)*Date{
	s := sd.clone()
	t := time.Date(s.Year+years, time.Month(s.Month+months), s.Day+days, s.Hour, s.Min, s.Sec, s.Nsec, s.Loc)
	return TimeToDate(t)
}

// addDay 以天数加,加days天
func (sd *Date)addDays(days int)*Date{
	return sd.clone().add(0,0, days)
}

// addMonths 以月数加,加months个月
func (sd *Date)addMonths(months int)*Date{
	return sd.clone().add(0,months, 0)
}

// addYears 以年数加,加years年
func (sd *Date)addYears(years int)*Date{
	return sd.clone().add(years,0, 0)
}

// nextYear 下一年
func (sd *Date)nextYear()*Date{
	return sd.clone().addYears(1)
}

// prevYear 前一年
func (sd *Date)prevYear()*Date{
	return sd.clone().addYears(-1)
}

// nextMonth 下一月
func (sd *Date)nextMonth()*Date{
	return sd.clone().addMonths(1)
}

// prevMonth 前一月
func (sd *Date)prevMonth()*Date{
	return sd.clone().addMonths(-1)
}

// nextDay 下一天
func (sd *Date)nextDay()*Date{
	return sd.clone().addDays(1)
}
// prevDay 前一天
func (sd *Date)prevDay()*Date{
	return sd.clone().addDays(-1)
}
// tomorrow 明天
func (sd *Date)tomorrow()*Date{
	return sd.clone().nextDay()
}
// yesterday 昨天
func (sd *Date)yesterday()*Date{
	return sd.clone().prevDay()
}

// monthFirstDayDate 该月1日
func (sd *Date)monthFirstDayDate()*Date{
	t := sd.clone().monthFirstDayTime()
	return TimeToDate(t)
}

// monthFirstDayDate 该月最后一日
func (sd *Date)monthLastDayDate()*Date{
	t := sd.clone().monthLastDayTime()
	return TimeToDate(t)
}

// monthDays 该月有多少天
func (sd *Date)monthDays()int{
	s := sd.clone()

	// 方法一,该月最后一天是几日就是该月天数
	// return s.MonthLastDay()

	// 方法二,该年是否为闰年,闰年则该年2月是29天,否则该年2月是28天
	md := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	k := s.Month - 1
	if s.Month == 2 && IsLeap(s.Year) {
		return 29
	}
	return md[k]
}

// monthFirstdayWeek 该月1日是星期几
func (sd *Date)monthFirstDayWeek()int{
	s := sd.clone()
	s.Day = 1
	w := s.time().Weekday()
	return int(w)
}

// monthLastdayWeek 该月最后一日是星期几
func (sd *Date)monthLastDayWeek()int{
	w := sd.clone().monthLastDayTime().Weekday()
	return int(w)
}

// monthLastDay 该月最后一天是几日
func (sd *Date)monthLastDay()int{
	return sd.clone().monthLastDayTime().Day()
}

// monthFirstDayTime 该月第一天 0:0:0.0 的Time
func (sd *Date)monthFirstDayTime()time.Time{
	s := sd.clone()
	s.Day = 1
	s.Hour = sd.Hour
	s.Min = sd.Min
	s.Sec = sd.Sec
	s.Nsec = sd.Nsec
	return s.time()
}

// monthLastDay 该月最后一天 23:59:59.999999999 的Time
func (sd *Date)monthLastDayTime()time.Time{
	s := sd.clone()
	return time.Date(s.Year, time.Month(s.Month+1), 1, 0, 0, 0, 0, s.Loc).Add(-1)
}

// time  Date转Time
func (sd *Date)time()time.Time{
	s := sd.clone()
	return time.Date(s.Year,time.Month(s.Month),s.Day,s.Hour,s.Min,s.Sec,s.Nsec,s.Loc)
}

// isoWeek 是该年第几周
func (sd *Date)isoWeek()(year, week int){
	return sd.clone().time().ISOWeek()
}