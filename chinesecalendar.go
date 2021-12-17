package gocalendar

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 中国(东八区)时间相对UTC的偏移量(单位：天days)
const cChineseTimeOffsetDays float64 = 8 / 24.0

// type GZItem 天干地支单元
type GZItem struct {
	HSI int    `json:"hsi"` // 天干索引
	HSN string `json:"hsn"` // 天干名称
	EBI int    `json:"ebi"` // 地支索引
	EBN string `json:"ebn"` // 地支名称
}

// type GZ 日期时间干支
type GZ struct {
	Year  *GZItem `json:"ygz"` // 年天干地支
	Month *GZItem `json:"mgz"` // 月天干地支
	Day   *GZItem `json:"dgz"` // 日天干地支
	Hour  *GZItem `json:"hgz"` // 时天干地支
}

// type LunarDate 农历
type LunarDate struct {
	Year          int      `json:"year"`      // 年
	Month         int      `json:"month"`     // 月
	Day           int      `json:"day"`       // 日
	MonthName     string   `json:"monthName"` // 月份名称
	DayName       string   `json:"dayName"`   // 日名称
	LeapStr       string   `json:"leapStr"`   // 闰字，可以用该值是否为空来判断该月是否为闰月
	YearLeapMonth int      `json:"ylm"`       // 该年闰几月，如果该年无闰月，则0，该年闰几月该值就是几，也可以用LeapMonth == Month判断该月是否闰月
	AnimalIndex   int      `json:"sai"`       // 年生肖索引
	AnimalName    string   `json:"san"`       // 年生肖名称
	YearGZ        *GZItem  `json:"ygz"`       // 年干支
	Festival      *FestivalItem `json:"festival"`  // 农历节日
}

// type pureJieQi16Temp struct pureJieSinceSpring和qiSinceWinterSolstice的缓存
type pureJieQi16Temp struct {
	data map[int][16]float64
	mu sync.RWMutex
}

// type trueNewMoon20Temp struct 20个新月点年表缓存
type trueNewMoon20Temp struct {
	data map[int][20]float64
	mu sync.RWMutex
}

// type lunarMonthCode15Temp struct农历月名称年表缓存
type lunarMonthCode15Temp struct {
	data map[int][15]float64
	mu sync.RWMutex
}

// type lunarMonthDays15Temp struct 农历月份对应的天数年表缓存
type lunarMonthDays15Temp struct {
	data map[int][15]int
	mu sync.RWMutex
}

// (*Calendar) ChineseSexagenaryCycle 日期时间对应的干支
//
// 特别提醒:干支推算的日干支存在早晚子时的区别
// NightZiHour默认为false是不区分早晚子时00:00-02:00为子时，NightZiHour为true时，23:00-24:00 00:00-01:00为子时
func (c *Calendar) ChineseSexagenaryCycle(t time.Time) GZ {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	// t时间时区与UTC的时差
	_,offset := t.Zone()
	var offsetDays = float64(offset) / 86400

	// t的儒略日
	// 该儒略日是按t的时区转换的
	jd := JulianDay(float64(year), float64(month), float64(day), float64(hour), float64(minute), float64(second))

	// 年月日时四柱
	// hsaeb := make(map[string]*GZItem)
	var gzs GZ

	// 立春点开始的节，年干支以立春开始(本方法为了日历使用，在计算中以立春当天为准，不考虑详细时间)
	// jss中儒略日是TT时间(这里强制为UTC时间)，是未经时区修改的儒略日
	// 在与jd比较时，应加上时区时差
	jss := c.pureJieSinceSpring(year)

	// 以立春当天0时作比较，不考虑定立春的时分秒，所以用math.Floor向下取整数部分
	if math.Floor(jd+0.5) < math.Floor(jss[1] + 0.5 + offsetDays) { // $jss[1]为立春，约在2月5日前后。
		year-- // 若小于jss[1]则属于前一个节气年

		// 取得自立春开始的节(不包含中气)，该数组长度固定为16
		jss = c.pureJieSinceSpring(year)
	}

	// 年干支
	ygz := ((year+4712+24)%60 + 60) % 60
	gzs.Year = &GZItem{
		HSI:  ygz % 10, // 年干
		EBI: ygz % 12, // 年支
	}

	ix := 0

	// 比较求算节气月，求出月干支
	for j := 0; j <= len(jss); j++ {
		if math.Floor(jss[j] + 0.5 + offsetDays) > math.Floor(jd+0.5) {
			// 已超过指定时刻，故应取前一个节气，用jd的0时比较，不考虑时分秒
			ix = j - 1
			break
		}
	}

	tmm := ((year + 4712) * 12 + (ix - 1) + 60) % 60 // 数组0为前一年的小寒所以这里再减一
	mgz := (tmm + 50) % 60
	gzs.Month = &GZItem{
		HSI:  mgz % 10, // 月干
		EBI: mgz % 12, // 月支
	}

	jdn := jd + 0.5                                  // 计算日柱的干支，加0.5是将起始点从正午改为0时开始
	thes := ((jdn - math.Floor(jdn)) * 86400) + 3600 // 将jd的小数部分化为秒，并加上起始点前移的一小时(3600秒)
	dayJd := math.Floor(jdn) + thes/86400            // 将秒数化为日数，加回到jd的整数部分
	dgz := (int(math.Floor(dayJd+49))%60 + 60) % 60
	gzs.Day = &GZItem{
		HSI:  dgz % 10, // 日干
		EBI: dgz % 12, // 日支
	}

	// 区分早晚子时,日柱前移一柱
	if c.config.NightZiHour && (hour >= 23) {
		gzs.Day = &GZItem{
			HSI:  (gzs.Day.HSI + 10 - 1) % 10,  // 日干
			EBI: (gzs.Day.EBI + 12 - 1) % 12, // 日支
		}
	}

	dh := (dayJd) * 12 // 计算时柱的干支
	hgz := (int(math.Floor(dh+48))%60 + 60) % 60
	gzs.Hour = &GZItem{
		HSI:  hgz % 10, // 时干
		EBI: hgz % 12, // 时支
	}

	// 为干支附中文名称
	gzs.Year.HSN   = heavenlyStemsNameArray[gzs.Year.HSI]
	gzs.Year.EBN  = earthlyBranchesNameArray[gzs.Year.EBI]
	gzs.Month.HSN  = heavenlyStemsNameArray[gzs.Month.HSI]
	gzs.Month.EBN = earthlyBranchesNameArray[gzs.Month.EBI]
	gzs.Day.HSN    = heavenlyStemsNameArray[gzs.Day.HSI]
	gzs.Day.EBN   = earthlyBranchesNameArray[gzs.Day.EBI]
	gzs.Hour.HSN   = heavenlyStemsNameArray[gzs.Hour.HSI]
	gzs.Hour.EBN  = earthlyBranchesNameArray[gzs.Hour.EBI]

	return gzs
}



// (*Calendar) pureJieSinceSpring 求出以某年立春点开始的节
func (c *Calendar) pureJieSinceSpring(year int) [16]float64 {
	// jss 16个节的jd数据

	// 如果c.jSS记录了该年的数据，则直接返回
	jss,err := c.tempData.jSS.getData(year)
	if err == nil {
		return jss
	}

	lastYearAsts := lastYearSolarTerms(float64(year))

	ki := -1 // 数组索引

	// 19小寒;21立春;23惊蛰
	for i := 19; i <= 23; i += 2 {
		// if i%2 == 0 {
		// 	continue
		// }
		if lastYearAsts[i] == 0 {
			continue
		}

		ki++
		// jss[ki] = Round(lastYearAsts[i]+cChineseTimeOffsetDays, 10) // 农历计算需要，加上中国(东八区)时差
		jss[ki] = lastYearAsts[i] // 中国(东八区)时差放在具体计算时调整，此处不做调整
	}

	asts := adjustedSolarTermsJd(float64(year), 0, 25)
	for i := 1; i <= 25; i += 2 {
		// if i%2 == 0 {
		// 	continue
		// }
		if asts[i] == 0 {
			continue
		}

		ki++
		// jss[ki] = Round(asts[i]+cChineseTimeOffsetDays, 10) // 农历计算需要，加上中国(东八区)时差
		jss[ki] = asts[i] // 中国(东八区)时差放在具体计算时调整，此处不做调整
	}

	c.tempData.jSS.setData(year,jss)

	return jss
}

// (*Calendar) qiSinceWinterSolstice 求出自上一年冬至点为起点的连续中气
func (c *Calendar) qiSinceWinterSolstice(year int) [16]float64 {
	// qss 16个中气的jd

	// 如果c.qSS记录了该年的数据，则直接返回
	qss,err := c.tempData.qSS.getData(year)
	if err == nil {
		return qss
	}

	lastYearAsts := lastYearSolarTerms(float64(year))

	ki := -1 // 数组索引

	// 18冬至(上一年);20大寒;22雨水
	for i := 18; i <= 22; i += 2 {
		if lastYearAsts[i] == 0 {
			continue
		}

		ki++
		// qss[ki] = Round(lastYearAsts[i]+cChineseTimeOffsetDays, 10) // 农历计算需要，加上中国(东八区)时差
		qss[ki] = lastYearAsts[i] // 中国(东八区)时差放在具体计算时调整，此处不做调整
	}

	asts := adjustedSolarTermsJd(float64(year), 0, 25)
	for i := 0; i <= 24; i += 2 {
		if asts[i] == 0 {
			continue
		}

		ki++
		// qss[ki] = Round(asts[i] + cChineseTimeOffsetDays, 10) // 农历计算需要，加上中国(东八区)时差
		qss[ki] = asts[i] // 中国(东八区)时差放在具体计算时调整，此处不做调整
	}

	c.tempData.qSS.setData(year,qss)

	return qss
}


// (*Calendar) GregorianToLunar 公历转农历
//
// @param int year  公历年份
// @param int month 公历月份
// @param int day   公历日
func (c *Calendar) GregorianToLunar(year, month, day int) LunarDate {
	t := time.Date(year,time.Month(month),day,0,0,0,0,c.loc)

	return c.gregorianToLunar(t,false) // festival:false时不取农历节日
}

// (*Calendar) gregorianToLunar 公历转农历
func (c *Calendar) gregorianToLunar(t time.Time,festival bool) LunarDate{
	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	lunarYear := year  // 初始农历年等于公历年

	prev := 0 // 是否跨年了,跨年了则减一
	isLeap := false // 是否闰月

	// t的儒略日
	jd := JulianDay(float64(year), float64(month), float64(day), float64(hour), float64(minute), float64(second))

	jdn := jd + 0.5 // 加0.5是将起始点从正午改为0时开始

	nm, lmc := c.zqAndSMandLunarMonthCode(year)

	// 如果公历日期的jd小于第一个朔望月新月点，表示农历年份是在公历年份的上一年
	if math.Floor(jdn) < math.Floor(nm[0] + 0.5 + cChineseTimeOffsetDays) {

		prev = 1
		nm, lmc = c.zqAndSMandLunarMonthCode(year-1)

	}

	// 查询对应的农历月份索引
	var mi = 0
	for i := 0; i <= 14; i++ { // 指令中加0.5是为了改为从0时算起而不从正午算起
		if math.Floor(jdn) >= math.Floor(nm[i]+0.5 + cChineseTimeOffsetDays) && math.Floor(jdn) < math.Floor(nm[i+1]+0.5 + cChineseTimeOffsetDays) {
			mi = i
			break
		}
	}

	// 农历的年
	// 如果月份属于上一年的11月或12月,或者农历年在上一年时
	if lmc[mi] < 2 || prev == 1 { // 年
		lunarYear --
	}

	// 农历月份是否是闰月
	if (lmc[mi] - math.Floor(lmc[mi])) * 2 + 1 != 1 { // 因mc(mi)=0对应到前一年农历11月,mc(mi)=1对应到前一年农历12月,mc(mi)=2对应到本年1月,依此类推
		isLeap = true
	}

	// 农历的月
	lunarMonth := int(math.Floor(lmc[mi] + 10)) % 12 + 1

	// 农历的日
	lunarDay := int(math.Floor(jdn) - math.Floor(nm[mi] + 0.5 + cChineseTimeOffsetDays) + 1) // 此处加1是因为每月初一从1开始而非从0开始


	// 整理年月日农历表示
	monthName := lunarMonthNameArray[lunarMonth - 1]
	dayName := DayChinese(lunarDay)
	ygz := ((lunarYear+4712+24)%60 + 60) % 60
	yhsi := ygz % 10
	yebi := ygz % 12
	animalIndex := yebi
	animalName := symbolicAnimalsNameArray[yebi]
	yearGZ := &GZItem{
		HSI:  yhsi,
		HSN:  heavenlyStemsNameArray[yhsi],
		EBI: yebi,
		EBN: earthlyBranchesNameArray[yebi],
	}

	// 整理闰月相关
	leapStr := "" //
	leapMonth := 0 // 闰几月
	if isLeap {
		leapStr = lunarLeapString
		leapMonth = lunarMonth
	}

	// 农历节日
	var lf FestivalItem
	if festival {
		lf = c.lunarFestival(lunarYear,lunarMonth,lunarDay,isLeap)
	}


	// 返回
	return LunarDate{
		Year:          lunarYear,
		Month:         lunarMonth,
		Day:           lunarDay,
		MonthName:     monthName,
		DayName:       dayName,
		LeapStr:       leapStr,
		YearLeapMonth: leapMonth,
		AnimalIndex:   animalIndex,
		AnimalName:    animalName,
		YearGZ:        yearGZ,
		Festival:      &lf,
	}
}


// (*Calendar) LunarToGregorian 农历转公历
//
// demo:
//     c := DefaultCalendar()
//     gd,_ := dc.LunarToGregorian(2020,4,14,true)
//     fmt.Println(gd.Format(time.RFC3339))
//
// @param int  lunarYear  农历年份
// @param int  lunarMonth 农历月份
// @param int  lunarDay   农历日
// @param bool isLeap   输入的日期是否是闰月的农历日期
func (c *Calendar) LunarToGregorian(lunarYear,lunarMonth,lunarDay int, isLeap bool) (time.Time, error){

	ld := LunarDate{Year: lunarYear, Month: lunarMonth,Day: lunarDay}
	if isLeap {
		ld.YearLeapMonth = lunarMonth
	}

	t,err := c.lunarToGregorian(ld)

	return t,err
}

// (*Calendar) lunarToGregorian 农历转公历
func (c *Calendar) lunarToGregorian(ld LunarDate) (time.Time, error){
	lunarYear,lunarMonth,lunarDay := ld.Year, ld.Month, ld.Day

	isLeap := false
	if ld.YearLeapMonth > 0 && ld.YearLeapMonth == ld.Month {
		isLeap = true
	}

	nm, lmc := c.zqAndSMandLunarMonthCode(lunarYear)

	// 该年闰几月，0无闰月
	leapMonth := mcLeap(lmc)

	// 11月对应到1,12月对应到2,1月对应到3,2月对应到4,依此类推
	lunarMonth += 2

	var nofd [15]int
	for i := 0; i <= 14; i++ {
		nofd[i] = int(math.Floor(nm[i+1] + 0.5 + cChineseTimeOffsetDays) - math.Floor(nm[i] + 0.5 + cChineseTimeOffsetDays)) // 每月天数,加0.5是因JD以正午起算
	}

	var jd float64

	if isLeap { // 闰月

		if leapMonth < 3 { // 而旗标非闰月或非本年闰月,则表示此年不含闰月.leap=0代表无闰月,=1代表闰月为前一年的11月,=2代表闰月为前一年的12月
			return time.Time{}, errors.New("此年非闰年") // 此年非闰年
		} else { // 若本年內有闰月
			if leapMonth != lunarMonth { // 但不为指入的月份
				return time.Time{}, errors.New("该月非闰月") // 则指定的月份非闰月,此月非闰月
			} else { // 若输入的月份即为闰月
				if lunarDay <= nofd[lunarMonth] { // 若指定的日期不大于当月的天數
					jd = nm[lunarMonth] + float64(lunarDay) - 1 // 则将当月之前的JD值加上日期之前的天數
				} else { // 日期超出范围
					return time.Time{}, errors.New("日期超出范围")
				}
			}
		}

	} else {
		if leapMonth == 0 { // 若旗标非闰月,则表示此年不含闰月(包括前一年的11月起之月份)
			if lunarDay <= nofd[lunarMonth-1] { // 若日期不大于当月天数
				jd = nm[lunarMonth-1] + float64(lunarDay) - 1 // 则将当月之前的JD值加上日期之前的天数
			} else { // 日期超出范围
				return time.Time{}, errors.New("日期超出范围")
			}
		} else { // 若旗标为本年有闰月(包括前一年的11月起之月份) 公式nofd(lunarMonth - (lunarMonth > leapMonth) - 1)的用意为:若指定月大于闰月,则索引用lunarMonth,否则索引用lunarMonth-1
			k := lunarMonth -1
			if lunarMonth > leapMonth {
				k = lunarMonth
			}
			if lunarDay <= nofd[k] { // 若输入的日期不大于当月天数
				jd = nm[k] + float64(lunarDay) - 1 // 则将当月之前的JD值加上日期之前的天数
			} else { // 日期超出范围
				return time.Time{}, errors.New("日期超出范围")
			}
		}
	}

	jd = math.Floor(jd) + 0.5
	return JdToTime(jd, c.loc), nil
}

// (*Calendar) LunarMonthDay 农历某个月有多少天
//
// @param int lunarYear   农历年数字
// @param int lunarMonth  农历月数字
// @param bool isLeap  是否是闰月
func (c *Calendar) LunarMonthDays(lunarYear,lunarMonth int, isLeap bool) (int,error) {

	var lmc [15]float64

	// 如果c.lMC记录了该年的数据，则直接赋值
	lMC,err := c.tempData.lMC.getData(lunarYear)
	if err == nil{
		lmc = lMC
	}else{
		_, lmc = c.zqAndSMandLunarMonthCode(lunarYear)
	}

	// 闰几月，0无闰月
	leapMonth := mcLeap(lmc)

	// 11月对应到1,12月对应到2,1月对应到3,2月对应到4,依此类推
	lunarMonth += 2

	lmd := c.mdList(lunarYear)

	dy := 0 // 当月天数

	if isLeap {
		if leapMonth < 3 { // 而旗标非闰月或非本年闰月,则表示此年不含闰月.leapMonth=0代表无闰月,=1代表闰月为前一年的11月,=2代表闰月为前一年的12月
			return 0, errors.New("该年非闰年")
		}
		// 若本年內有闰月
		if leapMonth != lunarMonth { // 但不为指定的月份
			return 0, errors.New("该月非该年的闰月")
		} else { // 若指定的月份即为闰月
			dy = lmd[lunarMonth]
		}
	} else { // 若没有指明是闰月
		k := lunarMonth -1
		if leapMonth != 0 {
			// 若旗标为本年有闰月(包括前一年的11月起之月份) 公式nofd(lunarMonth - (lunarMonth > leapMonth) - 1)的用意为:若指定月大于闰月,则索引用lunarMonth,否则索引用lunarMonth-1
			if lunarMonth > leapMonth {
				k = lunarMonth
			}
		}

		dy = lmd[k]
	}

	return dy, nil
}



// 农历一年的月份对应天数表
func (c *Calendar)mdList(lunarYear int) [15]int {
	// 如果c.lMD记录了该年的数据，则直接返回
	lmd,err := c.tempData.lMD.getData(lunarYear)
	if err == nil {
		return lmd
	}

	nm, _ := c.zqAndSMandLunarMonthCode(lunarYear)

	for i := 0; i <= 14; i++ {
		lmd[i] = int(math.Floor(nm[i+1] + 0.5 + cChineseTimeOffsetDays) - math.Floor(nm[i] + 0.5 + cChineseTimeOffsetDays)) // 每月天数,加0.5是因JD以正午起算
	}

	c.tempData.lMD.setData(lunarYear,lmd)

	return lmd
}

// (*Calendar) LunarLeap 取农历某年的闰月
//
// 0表示无闰月
func (c *Calendar) LunarLeap(lunarYear int) int {
	_,lmc := c.zqAndSMandLunarMonthCode(lunarYear)

	leap := mcLeap(lmc)

	return int(math.Max(0, float64(leap-2)))
}

// mcLeap 从农历的月代码lmc中找出闰月
//
// 0表示无闰月
func mcLeap (lmc [15]float64) int {

	var leap float64 = 0 // 若闰月旗标为0代表无闰月

	for j := 1; j <= 14; j++ { // 确认指定年前一年11月开始各月是否闰月
		if lmc[j]-math.Floor(lmc[j]) > 0 { // 若是,则将此闰月代码放入闰月旗标內
			leap = math.Floor(lmc[j] + 0.5) // leap = 0对应农历11月,1对应农历12月,2对应农历隔年1月,依此类推.
			break
		}
	}

	return int(leap)
}



// (*Calendar) zqAndSMandLunarMonthCode 以比较日期法求算冬月及其余各月名称代码,包含闰月,冬月为0,腊月为1,正月为2,其余类推.闰月多加0.5
//
// 农历按中国时间(东八区)计算
func (c *Calendar) zqAndSMandLunarMonthCode(year int) ([16]float64, [15]float64) {

	// 取得以前一年冬至为起点之连续16个中气
	qss := c.qiSinceWinterSolstice(year)

	// 求出以含冬至中气为阴历11月(冬月)开始的连续16个朔望月的新月点
	nm := c.sMsinceWinterSolstice(year, qss[0])

	// 如果c.lMC记录了该年的数据，则直接返回
	lmc,err := c.tempData.lMC.getData(year)
	if err == nil {
		return nm,lmc
	}

	// 设定旗标,0表示未遇到闰月,1表示已遇到闰月
	yz := 0

	if math.Floor(qss[12]+0.5 + cChineseTimeOffsetDays) >= math.Floor(nm[13]+0.5 + cChineseTimeOffsetDays) {

		for i := 1; i <= 14; i++ {

			// 至少有一个朔望月不含中气,第一个不含中气的月即为闰月
			// 若阴历腊月起始日大於冬至中气日,且阴历正月起始日小于或等于大寒中气日,则此月为闰月,其余同理
			if (nm[i]+0.5 + cChineseTimeOffsetDays) > math.Floor(qss[i-1-yz]+0.5 + cChineseTimeOffsetDays) && math.Floor(nm[i+1]+0.5 + cChineseTimeOffsetDays) <= math.Floor(qss[i-yz]+0.5 + cChineseTimeOffsetDays) {
				lmc[i] = float64(i) - 0.5
				yz = 1 // 标示遇到闰月
			} else {
				lmc[i] = float64(i - yz) // 遇到闰月开始,每个月号要减1
			}
		}
	} else { // 否则表示两个连续冬至之间只有11个整月,故无闰月

		for i := 0; i <= 12; i++ { // 直接赋予这12个月月代码
			lmc[i] = float64(i)
		}
		for i := 13; i <= 14; i++ { // 处理次一置月年的11月与12月,亦有可能含闰月
			// 若次一阴历腊月起始日大于附近的冬至中气日,且阴历正月起始日小于或等于大寒中气日,则此月为腊月,次一正月同理.
			if (nm[i]+0.5 + cChineseTimeOffsetDays) > math.Floor(qss[i-1-yz]+0.5 + cChineseTimeOffsetDays) && math.Floor(nm[i+1]+0.5 + cChineseTimeOffsetDays) <= math.Floor(qss[i-yz]+0.5 + cChineseTimeOffsetDays) {
				lmc[i] = float64(i) - 0.5
				yz = 1 // 标示遇到闰月
			} else {
				lmc[i] = float64(i - yz) // 遇到闰月开始,每个月号要减1
			}
		}
	}

	c.tempData.lMC.setData(year,lmc)

	return nm, lmc
}



// (*Calendar) sMsinceWinterSolstice 求算以含冬至中气为阴历11月开始的连续16个朔望月
func (c *Calendar) sMsinceWinterSolstice (year int, dzJd float64) [16]float64 {

	tnm := [20]float64{}
	nm := [16]float64{}

	// 如果c.tNM记录了该年的数据，则直接赋值给tnm
	tNM,err := c.tempData.tNM.getData(year)
	if err == nil{
		tnm = tNM
	}else{

		// 求年初前两个月附近的新月点(即前一年的11月初)
		novemberJd := JulianDay(float64(year) - 1, 11, 1)

		// 求得自2000年1月起第kn个平均朔望日及其JD值
		// kn,thejd := meanNewMoon(novemberJd)
		kn := referenceLunarMonthNum(novemberJd)

		// 求出连续20个朔望月
		for i := 0; i <= 19; i++ {
			k := kn + float64(i)

			// 以k值代入求瞬时朔望日
			// tnm[i] = trueNewMoon(k) + cChineseTimeOffsetDays // 农历计算需要，加上中国(东八区)时差
			tnm[i] = trueNewMoon(k) // 中国(东八区)时差放在具体计算时调整，此处不做调整

			// 下式为修正 dynamical time to Universal time
			// 1为1月，0为前一年12月，-1为前一年11月(当i=0时，i-1代表前一年11月)
			tnm[i] = Round(tnm[i] - deltaTDays(float64(year), float64(i - 1)), 10)
		}

		c.tempData.tNM.setData(year,tnm)
	}

	var jj = 0
	for j := 0; j <= 18; j++ {
		if math.Floor(tnm[j] + 0.5) > math.Floor(dzJd + 0.5) {
			jj = j
			break
		} // 已超过冬至中气(比较日期法)
	}

	for k := 0; k <= 15; k++ { // 取上一步的索引值
		nm[k] = tnm[jj-1+k] // 重排索引,使含冬至朔望月的索引为0
	}

	return nm
}

// (*Calendar) lunarFestival 取农历节日
func (c *Calendar) lunarFestival (lunarYear,lunarMonth,lunarDay int, isLeap bool) FestivalItem {

	fds := c.tempData.lFD.getData(lunarYear)

	if len(fds) == 0 {

		// 根据lunarFestivalArray重新格式一个准确的月日为索引的节日map
		for lfK,lfV := range lunarFestivalArray{
			trueK := "" // 经过处理转换成月日的K
			isleap := false    // 索引中是否指明为闰月
			isLastDay := false // 索此中是否指明为某月最后一天
			m := ""  // 月
			d := ""  // 日

			// 索引正则
			// re := regexp.MustCompile("([0-9]{1,2})?(@?)(M)?(?:([0-9]{1,2})D)?(?:([1-4])W([0-6]))?(\\$)?$").FindStringSubmatch(lfK)
			// 农历节日索引正则
			re := regexp.MustCompile("([0-9]{1,2})(@?)M(?:([0-9]{1,2})D)?(\\$)?$").FindStringSubmatch(lfK)

			if len(re) != 5{
				continue // 索引格式不正确
			}
			// 月数
			if re[1] == "" {
				continue // 索引格式不正确,没指明月份
			}
			m = re[1]
			// 月数string转int
			month,err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				continue // 月份为空或数字不正确
			}

			// 开始拼接trueK
			trueK = m

			// 闰月
			if re[2] != "" {
				isleap = true

				lm := c.LunarLeap(lunarYear)
				if lm == 0 && lm != month {
					continue // 当前年中无该闰月,不用记录该节日
				}

				trueK += "@"
			}

			trueK += "M"

			// 日数
			if re[3] != "" {
				d = re[3]
			}
			// 最后一日
			if re[4] != "" && re[3] == "" {
				isLastDay = true
			}

			// 如果没指明日数也没有指明是最后一天，则索引是无效的
			if d == "" && !isLastDay {
				continue
			}

			// m月有多少天，验证d是否大于这个值以及将用该值表示最后一天
			days,err := c.LunarMonthDays(lunarYear,month,isleap)
			if err != nil {
				continue
			}

			day := 0
			if d != "" {
				day,err = strconv.Atoi(d)
				if err != nil || day < 1 || day > days {
					continue // 日为空或数字不正确
				}
			}else {
				day = days // 最后一日
			}

			trueK += strconv.Itoa(day) + "D"

			// 对应值
			if _, ok := fds[trueK]; ok {
				fds[trueK] = append(fds[trueK], strings.Split(lfV,",")...)
			}else{
				fds[trueK] = strings.Split(lfV,",")
			}
		}

		c.tempData.lFD.setData(lunarYear,fds)
	}

	festivalIndex := strconv.Itoa(lunarMonth)
	if isLeap {
		festivalIndex += "@"
	}

	festivalIndex += "M" + strconv.Itoa(lunarDay) + "D"


	var fi FestivalItem // 该日的节日
	if fv, ok := fds[festivalIndex]; ok {
		for _, v := range fv {
			v = strings.TrimSpace(v)
			if svs := strings.Split(v,"*"); len(svs) > 1 {
				fi.Show = append(fi.Show,svs[1])
			}else{
				fi.Secondary = append(fi.Secondary,v)
			}
		}
	}

	return fi
}

// DayChinese 农历日汉字表示法
func DayChinese(d int) string {
	daystr := ""
	if d < 1 || d > 30 {
		return ""
	}
	switch d {
	case 10:
		daystr = lunarWholeTensArray[0] + lunarNumberArray[10]
	case 20:
		daystr = lunarNumberArray[2] + lunarNumberArray[10]
	case 30:
		daystr = lunarNumberArray[3] + lunarNumberArray[10]
	default:
		k := d / 10
		m := d % 10
		daystr = lunarWholeTensArray[k] + lunarNumberArray[m]
	}
	return daystr
}

// (GZ) String 干支显示
func (gz GZ) String() string{
	return fmt.Sprintf("%s%s年%s%s月%s%s日%s%s时", gz.Year.HSN, gz.Year.EBN, gz.Month.HSN, gz.Month.EBN, gz.Day.HSN, gz.Day.EBN, gz.Hour.HSN, gz.Hour.EBN)
}

// (LunarDate) String 农历显示
func (ld LunarDate)String()string{
	festivalStr := ""
	if ld.Festival != nil && ld.Festival.Show != nil{
		festivalStr = " " + strings.Join(ld.Festival.Show, ",")
	}

	return fmt.Sprintf("%d%s%s(%s)年%s%s月%s%s",ld.Year, ld.YearGZ.HSN, ld.YearGZ.EBN, ld.AnimalName, ld.LeapStr, ld.MonthName, ld.DayName, festivalStr)
}




// (*pureJieQi16) getData 读节气缓存年表
func (pjq *pureJieQi16Temp) getData(k int) ([16]float64, error) {
	pjq.mu.RLock()
	defer pjq.mu.RUnlock()
	var rv [16]float64
	if pjq.data == nil {
		pjq.data = make(map[int][16]float64)
	}
	if v, ok := pjq.data[k]; ok {
		return v,nil
	}

	return rv,errors.New("缓存数据不存在！")
}

// (*pureJieQi16) getData 写节气缓存年表
func (pjq *pureJieQi16Temp) setData (k int, v [16]float64){
	pjq.mu.Lock()
	defer pjq.mu.Unlock()
	if pjq.data == nil {
		pjq.data = make(map[int][16]float64)
	}
	pjq.data[k] = v
}

// (*trueNewMoon20Temp) getData 读20个新月点缓存年表
func (tnm *trueNewMoon20Temp) getData(k int) ([20]float64, error) {
	tnm.mu.RLock()
	defer tnm.mu.RUnlock()
	var rv [20]float64
	if tnm.data == nil {
		tnm.data = make(map[int][20]float64)
	}
	if v, ok := tnm.data[k]; ok {
		return v,nil
	}

	return rv,errors.New("缓存数据不存在！")
}

// (*trueNewMoon20Temp) getData 写20个新月点缓存年表
func (tnm *trueNewMoon20Temp) setData (k int, v [20]float64){
	tnm.mu.Lock()
	defer tnm.mu.Unlock()
	if tnm.data == nil {
		tnm.data = make(map[int][20]float64)
	}
	tnm.data[k] = v
}

// (*lunarMonthCode15Temp) getData 读农历月份代码缓存年表
func (mc *lunarMonthCode15Temp) getData(k int) ([15]float64, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	var rv [15]float64
	if mc.data == nil {
		mc.data = make(map[int][15]float64)
	}
	if v, ok := mc.data[k]; ok {
		return v,nil
	}

	return rv,errors.New("缓存数据不存在！")
}

// (*lunarMonthCode15Temp) getData 写农历月份代码缓存年表
func (mc *lunarMonthCode15Temp) setData (k int, v [15]float64){
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if mc.data == nil {
		mc.data = make(map[int][15]float64)
	}
	mc.data[k] = v
}

//
// (*lunarMonthDays15Temp) getData 读农历月份天数缓存年表
func (md *lunarMonthDays15Temp) getData(k int) ([15]int, error) {
	md.mu.RLock()
	defer md.mu.RUnlock()
	var rv [15]int
	if md.data == nil {
		md.data = make(map[int][15]int)
	}
	if v, ok := md.data[k]; ok {
		return v,nil
	}

	return rv,errors.New("缓存数据不存在！")
}

// (*lunarMonthDays15Temp) getData 写农历月份天数缓存年表
func (md *lunarMonthDays15Temp) setData (k int, v [15]int){
	md.mu.Lock()
	defer md.mu.Unlock()
	if md.data == nil {
		md.data = make(map[int][15]int)
	}
	md.data[k] = v
}