package gocalendar

import (
	"time"
)

type Ctime time.Time


func NewCtime() Ctime {
	return Ctime(time.Now())
}

func DateNewCtime(da *Date) Ctime {

	if da.Year == 0 {
		da.Year = 1970
	}
	if da.Month == 0 {
		da.Month = 1
	}
	if da.Day == 0 {
		da.Day = 1
	}
	if da.Loc == nil {
		da.Loc = time.Local
	}

	return Ctime(time.Date(da.Year, time.Month(da.Month), da.Day, da.Hour, da.Min, da.Sec, da.Nsec, da.Loc))
}

func UnixNewCtime(sec int64, nsec int64) Ctime {
	return Ctime(time.Unix(sec,nsec))
}

func (ct Ctime)Location()*time.Location{
	return time.Time(ct).Location()
}

func (ct Ctime) UTC() Ctime {
	return Ctime(time.Time(ct).UTC())
}
func (ct Ctime) Local() Ctime {
	return Ctime(time.Time(ct).Local())
}

func (ct Ctime) In(loc *time.Location) Ctime {
	return Ctime(time.Time(ct).In(loc))
}

func (ct Ctime) Zone() (name string, offset int) {
	return time.Time(ct).Zone()
}

// unix时间戳
func (ct Ctime) Unix() int64 {
	return time.Time(ct).Unix()
}

func (ct Ctime) Format(layout string) string {
	return time.Time(ct).Format(layout)
}



// 公历年是否为闰年
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// golang 1000-1-1之前的年份都归为1
func (ct Ctime) Year() int {
	return time.Time(ct).Year()
}

func (ct Ctime) Month() int {
	return int(time.Time(ct).Month())
}

// 一个月有多少天
// golang 对1582年10月4日至1582年10月15日少的10天加进去了
// 1582-10-04 23:59:59的unix时间戳是-12220156801
// 1582-10-15 00:00:00的unix时间戳是-12219292800
// 相差10天
func (ct Ctime) MonthDay() int {
	year, month, _ := time.Time(ct).Date()
	m := int(month)

	// 每月多少天的数组，索引为m-1
	md := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	d := md[m-1]

	if m == 2 && IsLeapYear(year) {
		d += 1
	}

	return d

	/*if year == 1582 && m == 10 {
		return 21
	}

	ndf1 := -B2i(year % 4 == 0)
	ndf2 := B2i(I2b(B2i(year % 400 == 0) - B2i(year % 100 == 0)) && (year > 1582))
	ndf := ndf1 + ndf2

	return 30 + (int(math.Abs(float64(m) - 7.5)+ 0.5) % 2) - B2i(m == 2) * (2 + ndf)*/
}

// 所在月还剩几天，不包括所在天在内
func (ct Ctime) MonthDayLeft() int {
	return ct.MonthDay() - ct.Day()
}

func (ct Ctime) Date() (year, month, day int) {
	year, m, day := time.Time(ct).Date()
	month = int(m)
	return year, month, day
}

func (ct Ctime) Day() int {
	return time.Time(ct).Day()
}

func (ct Ctime) Hour() int {
	return time.Time(ct).Hour()
}

func (ct Ctime) Second() int {
	return time.Time(ct).Second()
}

func (ct Ctime) Minute() int {
	return time.Time(ct).Minute()
}

func (ct Ctime) Nanosecond() int {
	return time.Time(ct).Nanosecond()
}

func (ct Ctime) Clock() (hour, min, sec int) {
	return time.Time(ct).Clock()
}

// 所在日是所在年的第几天
func (ct Ctime) YearDay() int {
	return time.Time(ct).YearDay()
}

// 所在周是所在年的第几周
func (ct Ctime) ISOWeek() (year, week int) {
	return time.Time(ct).ISOWeek()
}

// 星期几
func (ct Ctime) Weekday() int {
	return int(time.Time(ct).Weekday())
}

func (ct Ctime) AddDate(years, months, days int) Ctime {
	return Ctime(time.Time(ct).AddDate(years, months, days))
}

func (ct Ctime) PrevYear() (year int) {
	year, _, _ = ct.AddDate(-1, 0, 0).Date()
	return year
}

func (ct Ctime) NextYear() (year int) {
	year, _, _ = ct.AddDate(1, 0, 0).Date()
	return year
}

func (ct Ctime) PrevMonth() (year, month int) {
	year, month, _ = ct.AddDate(0, -1, 0).Date()
	return year, month
}

func (ct Ctime) NextMonth() (year, month int) {
	year, month, _ = ct.AddDate(0, 1, 0).Date()
	return year, month
}

func (ct Ctime) PrevDay() (year, month, day int) {
	year, month, day = ct.AddDate(0, 0, -1).Date()
	return year, month, day
}

func (ct Ctime) NextDay() (year, month, day int) {
	year, month, day = ct.AddDate(0, 0, 1).Date()
	return year, month, day
}

func B2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
func I2b(i int) bool {
	if i == 0 {
		return false
	}
	return true
}
