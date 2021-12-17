package gocalendar

import (
	"math"
	"time"
)

const (

	// 儒略历历法废弃年
	cJulianAbandonmentYear float64 = 1582

	// 儒略历历法废弃月
	cJulianAbandonmentMonth float64 = 10

	// 儒略历历法废弃日
	cJulianAbandonmentDay float64 = 4

	// 格里历历法实施年
	cGregorianAdoptionYear float64 = 1582

	// 格里历历法实施月
	cGregorianAdoptionMonth float64 = 10

	// 格里历历法实施日
	cGregorianAdoptionDay float64 = 15

	// 格里历历法实施日期TT时间(1582年10月15日)中午12点的儒略日
	cJulianGregorianBoundary float64 = 2299161.0

	// J2000.0的儒略日 (TT时间2000年1月1日中午12点 (UTC时间2000年1月1日11:58:55.816)的儒略日)
	cJulianDayJ2000 = 2451545.0

	// 儒略历1年有多少天
	cDaysOfAYear = 365.25
)

// JulianDay 计算日期时间(TT)的儒略日
//
// (特别提醒,我们将一个日期时间转为儒略日时,其实使用的并不是真正的TT时间，而是我们常用的UTC或当地时区时间,故此无需考虑TT与UTC之间的转换)
func JulianDay(year, month, day float64, timeParts ...float64) float64 {

	var hour float64 = 0
	var minute float64 = 0
	var second float64 = 0
	var millisecond float64 = 0

	for timeIndex, timePart := range timeParts {
		switch timeIndex {
		case 0:
			hour = timePart
		case 1:
			minute = timePart
		case 2:
			second = timePart
		case 3:
			millisecond = timePart
		}
	}

	// 计算公式参见: https://zh.wikipedia.org/wiki/%E5%84%92%E7%95%A5%E6%97%A5
	// 或参见: https://blog.csdn.net/weixin_42763614/article/details/82880007

	a := math.Floor((14 - month) / 12)
	y := year + 4800 - a
	m := month + 12 * a - 3
	second += millisecond / 1000.0
	d := day + hour / 24.0 + minute / 1440.0 + second / 86400.0

	var jdn float64

	// 依据儒略历废弃日期和格里历实施日期，使用两个不同的公式计算儒略日JDN
	if year < cJulianAbandonmentYear || (year == cJulianAbandonmentYear && month < cJulianAbandonmentMonth) || (year == cJulianAbandonmentYear && month == cJulianAbandonmentMonth && day <= cJulianAbandonmentDay) {
		// 儒略历日期
		jdn = jdnInJulian(y, m, d)
	} else if year > cGregorianAdoptionYear || (year == cGregorianAdoptionYear && month > cGregorianAdoptionMonth) || (year == cGregorianAdoptionYear && month == cGregorianAdoptionMonth && day >= cGregorianAdoptionDay) {
		// 格里历日期
		jdn = jdnInGregorian(y, m, d)
	} else {
		// 在儒略历废弃与格里历实施这中间有一段日期，这段日期的儒略日计算使用格里历实施的起始日计算
		jdn = jdnInGregorian(cGregorianAdoptionYear, cGregorianAdoptionMonth, cGregorianAdoptionDay)
	}

	return Round(jdn - 0.5, 10)
}

// jdnInJulian 儒略历日期(TT)转儒略日JDN
//
// JDN表达式与JD的关系是: JDN = JD + 0.5
func jdnInJulian(year, month, day float64) float64 {

	return day + math.Floor((153 * month + 2) / 5) + 365 * year + math.Floor(year / 4) - 32083
}

// jdnInGregorian 格里历日期(TT)转儒略日JDN
//
// JDN表达式与JD的关系是: JDN = JD + 0.5
func jdnInGregorian(year, month, day float64) float64 {

	return day + math.Floor((153 * month + 2) / 5) + 365 * year + math.Floor(year / 4) - math.Floor(year / 100) + math.Floor(year / 400) - 32045
}

// JdToTimeMap 儒略日计算对应的日期时间(TT)
//
// (特别提醒,我们将一个日期时间转为儒略日时,其实使用的并不是真正的TT时间，而是我们常用的UTC或当地时区时间,故此无需考虑TT与UTC之间的转换)
func JdToTimeMap(jd float64) map[string]int {

	jdn := jd + 0.5

	// 计算公式: https://blog.csdn.net/weixin_42763614/article/details/82880007
	Z := math.Floor(jdn) // 儒略日的整数部分
	F := jdn - Z         // 儒略日的小数部分

	var A float64

	// 2299161 是1582年10月15日12时0分0秒
	if Z < cJulianGregorianBoundary {
		// 儒略历
		A = Z
	} else {
		a := math.Floor((Z - 2305507.25) / 36524.25)
		var IntervalDays float64 = 10 // 儒略历被废弃至格里历被启用之间相差多少天
		A = Z + IntervalDays + a - math.Floor(a/4)
	}

	var dayF float64 = 1
	var C float64 = 0
	var E float64 = 0
	var k float64 = 0

	for {
		B := A + 1524                              // 以BC4717年3月1日0时为历元
		C = math.Floor((B - 122.1) / cDaysOfAYear) // 积年
		D := math.Floor(cDaysOfAYear * C)          // 积年的日数
		E = math.Floor((B - D) / 30.6)             // B-D为年内积日，E即月数
		dayF = B - D - math.Floor(30.6*E) + F

		// 否则即在上一月，可前置一日重新计算
		if dayF >= 1 {
			break
		}

		A -= 1
		k += 1
	}

	var year float64  // 年
	var month float64 // 月

	if E < 14 {
		month = E - 1
	} else {
		month = E - 13
	}

	if month > 2 {
		year = C - 4716
	} else {
		year = C - 4715
	}

	dayF += k
	if dayF == 0 {
		dayF += 1
	}

	// 天数分开成天与时分秒
	day := math.Floor(dayF) // 天
	dayD := dayF - day

	var hh, ii, ss, ms float64

	if dayD > 0 {
		hhF := dayD * 24 + 0.000000005 // 0.000000005补一个精度差
		hh = math.Floor(hhF)         // 时
		hhD := hhF - hh
		if hhD > 0 {
			iiF := hhD * 60
			ii = math.Floor(iiF) // 分
			iiD := iiF - ii
			if iiD > 0 {
				ssF := iiD * 60
				ss = math.Floor(ssF) // 秒
				ssD := ssF - ss
				if ssD > 0 {
					ms = ssD * 1000 // 毫秒
				}
			}
		}
	}

	return map[string]int{
		"year":        int(year),
		"month":       int(month),
		"day":         int(day),
		"hour":        int(hh),
		"minute":      int(ii),
		"second":      int(ss),
		"millisecond": int(ms),
	}
}

// TimeMapToTime 将time map转成loc时区的 *time.Time
//
// 该方法未将TT转为UTC，而是将TT强行等于UTC，如需TT日期，请不要转换,
// (TT与UTC的计算公式可以参考 TT = UTC + 64.184s ,
// 但由于我们在将日期时间转换为儒略日时使用的并不是真正的TT时间，而是UTC或当地时区时间，故此处不做TT与UTC转换.)
func TimeMapToTime(tm map[string]int, loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.Date(tm["year"],time.Month(tm["month"]),tm["day"],tm["hour"],tm["minute"],tm["second"],0,time.UTC).In(loc)
}

// JdToTime 将儒略日转成loc时区的 *time.Time
//
// 该方法未将TT转为UTC，而是将TT等于UTC，如需TT日期，请使用 JdToTimeMap 方法
func JdToTime(jd float64, loc *time.Location) time.Time {
	tm := JdToTimeMap(jd)
	return TimeMapToTime(tm,loc)
}

// Mjd 计算日期时间(TT)的简化儒略日
//
// 简化儒略日(Modified Julian Day, MJD)是将儒略日(Julian Day, JD)进行简化后得到的新计时法。
// 1957年,简化儒略日由史密松天体物理台(Smithsonian Astrophysical Observatory)引入。
// 1957年史密松天体物理台为便用于记录“伴侣号”人造卫星的轨道,将儒略日进行了简化，并将其命名为简化儒略日,其定义为: MJD=JD-2400000.5
// 儒略日2400000是1858年11月16日,因为JD从中午开始计算,所以简化儒略日的定义中引入偏移量0.5,
// 这意味着MJD 0相当于1858年11月17日的凌晨,并且每一个简化儒略日都在世界时午夜开始和结束.
func Mjd(year, month, day float64, timeParts ...float64) float64 {

	jd := JulianDay(year, month, day, timeParts...)

	return jdToMjd(jd)
}

// MjdToTimeMap 简化儒略日计算对应的日期时间(TT)
func MjdToTimeMap(mjd float64) map[string]int {

	jd := mjdToJd(mjd)
	return JdToTimeMap(jd)
}

// jdToMjd 儒略日转换为简儒略日
func jdToMjd(jd float64) float64 {
	return Round(jd - 2400000.5,10)
}

// mjdToJd 简儒略日转换为儒略日
func mjdToJd(mjd float64) float64 {
	return Round(mjd + 2400000.5,10)
}

// julianCentury 计算标准历元起的儒略世纪
func julianCentury(jd float64) float64 {
	return julianDayFromJ2000(jd) / cDaysOfAYear / 100.0
}

// julianThousandYear 计算标准历元起的儒略千年数
func julianThousandYear(jd float64) float64 {
	return julianDayFromJ2000(jd) / cDaysOfAYear / 1000.0
}

// julianDayFromJ2000 计算标准历元起的儒略日
func julianDayFromJ2000(jd float64) float64 {
	return jd - cJulianDayJ2000
}
