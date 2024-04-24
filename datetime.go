package gocalendar

import (
	"time"
)

// GregorianMonthDays 公历某月的总天数
func GregorianMonthDays(year, month int) int {
	md := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if month == 2 && IsLeapYear(year) {
		return 29
	}
	return md[month - 1]
}

// IsLeapYear 给定的公历年year是否是闰年

// 增加3200年的情况
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0) && (year%3200 != 0 || year%172800 == 0)
}

// BeginningOfMonth 计算出t所在月份的首日
//
// 时hour,分minute,秒second与t对应的值一样
func BeginningOfMonth(t time.Time) time.Time {
	day := t.Day()
	day--
	return t.AddDate(0,0,-day)
}

// EndOfMonth 计算出t所在月份的最后一日
//
// 时hour,分minute,秒second与t对应的值一样
func EndOfMonth(t time.Time) time.Time {
	nextMonthTime := t.AddDate(0,1,0)
	return BeginningOfMonth(nextMonthTime).AddDate(0,0,-1)
}