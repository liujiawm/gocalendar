/**
一个用golang写的日历，有公历转农历，农历转公历，节气，干支，星座，生肖等功能
中国的农历历法综合了太阳历和月亮历,为中国的生活生产提供了重要的帮助,是中国古人智慧与中国传统文化的一个重要体现

程序比较准确的计算出农历与二十四节气(精确到分),时间限制在1000-3000年间,在实际使用中注意限制年份
*/

package gocalendar

import (
	"errors"
	"fmt"
	"math"
	"time"
)

const Author = "liujiawm@gmail.com"
const Version = "1.0.2"

// 取日历方式
const (
	GridMonth int = iota
	GridWeek
	GridDay
)

// 日历
type Calendar struct {
	Loc       *time.Location // time.Location 默认time.Local
	FirstWeek int            // 日历显示时第一列显示周几，(日历表第一列是周几,0周日,依次最大值6)
	Grid      int            // 取日历方式,GridDay按天取日历,GridWeek按周取日历,GridMonth按月取日历
	Zwz       bool           // 是否区分早晚子时(子时从23-01时),true则23:00-24:00算成上一天
	Getjq     bool
}

type CalendarData struct {
	SD *SolarDate `json:"solar"`
	LD *LunarDate `json:"lunar"`
}

// 公历
type SolarDate struct {
	*Date  `json:"date"`
	Jq     *SolarJQ           `json:"jq"`
	GanZhi *SolarTianGanDiZhi `json:"gan_zhi"`
}

// 节气及时间
type SolarJQ struct {
	Name string `json:"name"`
	Date *Date  `json:"date"`
}

// 公历天干地支
type SolarTianGanDiZhi struct {
	Ytg    int    `json:"ytg"`     // 年天干
	YtgStr string `json:"ytg_str"` // 年天干名称
	Ydz    int    `json:"ydz"`     // 年地支
	YdzStr string `json:"ydz_str"` // 年地支名称
	Mtg    int    `json:"mtg"`     // 月天干
	MtgStr string `json:"mtg_str"` // 月天干名称
	Mdz    int    `json:"mdz"`     // 月地支
	MdzStr string `json:"mdz_str"` // 月地支名称
	Dtg    int    `json:"dtg"`     // 日天干
	DtgStr string `json:"dtg_str"` // 日天干名称
	Ddz    int    `json:"ddz"`     // 日地支
	DdzStr string `json:"ddz_str"` // 日地支名称
	Htg    int    `json:"htg"`     // 时天干
	HtgStr string `json:"htg_str"` // 时天干名称
	Hdz    int    `json:"hdz"`     // 时地支
	HdzStr string `json:"hdz_str"` // 时地支名称
}

// 农历
type LunarDate struct {
	*Date     `json:"date"`
	MonthStr  string          `json:"month_str"`  // 月的农历名称
	DayStr    string          `json:"day_str"`    // 天的农历名称
	LeapStr   string          `json:"leap_str"`   // 闰
	MonthDays int             `json:"-"`          // 当月有多少天
	LeapYear  int             `json:"leap_year"`  // 是否闰年，0不是闰年，大于就是闰几月
	LeapMonth int             `json:"leap_month"` // 当前前是否是所闰的那个月，0不是，1本月就是闰月
	YearGanZi *LunarYearGanZi `json:"gan_zi"`
}

// 农历通谷记年(干支和生肖属相)
type LunarYearGanZi struct {
	Gan    string `json:"gan"`
	Zhi    string `json:"zhi"`
	Animal string `json:"animal"`
}

var (
	LeapStr = "闰"
	// 中文数字
	NumberChineseArray = [11]string{"日", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}

	// 农历月份常用称呼
	MonthChineseArray = [12]string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "冬", "腊"}

	// 农历日期常用称呼
	DayChineseArray = [4]string{"初", "十", "廿", "卅"}

	// 天干
	TianGanArray = [10]string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

	// 地支
	DiZhiArray = [12]string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 生肖
	SymbolicAnimalsArray = [12]string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}

	// 节气
	JieQiArray = [24]string{"春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露",
		"秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至", "小寒", "大寒", "立春", "雨水", "惊蛰"}

	// 均值朔望月长(mean length of synodic month)
	// 朔望月长每个月都不同，此处仅为均值，在算出新月点后，还需加上一个调整值
	synMonth float64 = 29.530588853

	// 若以2000年的第一个均值新月点为基准点，此基准点为 2000年1月6日14时20分37秒(转换为中国时间是2000年1月6日22时20分37秒)。
	// 其对应的真实新月点为2000年1月6日18时13分42秒(转换为中国时间是2000年1月7日2时13分42秒)。
	// 此作为基准点的均值新月点的JDE值为bnm=2451550.09765日,定此为第0個新月点。
	bnm float64 = 2451550.09765

	// 因子
	ptsa = [...]float64{485, 203, 199, 182, 156, 136, 77, 74, 70, 58, 52, 50, 45, 44, 29, 18, 17, 16, 14, 12, 12, 12, 9, 8}
	ptsb = [...]float64{324.96, 337.23, 342.08, 27.85, 73.14, 171.52, 222.54, 296.72, 243.58, 119.81, 297.17, 21.02, 247.54,
		325.15, 60.93, 155.12, 288.79, 198.04, 199.76, 95.39, 287.11, 320.81, 227.73, 15.45}
	ptsc = [...]float64{1934.136, 32964.467, 20.186, 445267.112, 45036.886, 22518.443, 65928.934, 3034.906, 9037.513, 33718.147,
		150.678, 2281.226, 29929.562, 31555.956, 4443.417, 67555.328, 4562.452, 62894.029, 31436.921, 14577.848, 31931.756,
		34777.259, 1222.114, 16859.074}
)

// DefaultCalendar 默认日历设置
func DefaultCalendar() *Calendar {
	dc := &Calendar{
		Loc:       time.Local,
		FirstWeek: 0,
		Grid:      GridMonth,
		Zwz:       false,
		Getjq:     true,
	}
	return NewCalendar(dc)
}

// NewCalendar 日历设置
func NewCalendar(c *Calendar) *Calendar {
	if c.Loc == nil {
		c.Loc = time.Local
	}
	if c.FirstWeek < 0 || c.FirstWeek > 6 {
		c.FirstWeek = 0
	}
	return &Calendar{
		Loc:       c.Loc,
		FirstWeek: c.FirstWeek,
		Grid:      c.Grid,
		Zwz:       c.Zwz,
		Getjq:     c.Getjq,
	}
}

// NowCalendars 当前时间的日历,包括公历和农历
func (sc *Calendar) NowCalendars() []*CalendarData {
	if sc.Loc == nil {
		sc.Loc = time.Local
	}
	if sc.FirstWeek < 0 || sc.FirstWeek > 6 {
		sc.FirstWeek = 0
	}

	nt := timeFun().In(sc.Loc)
	return sc.calendars(TimeToDate(nt))
}

// Calendars 日历,包括公历和农历
func (sc *Calendar) Calendars(y, m, d int) []*CalendarData {
	if sc.Loc == nil {
		sc.Loc = time.Local
	}
	if sc.FirstWeek < 0 || sc.FirstWeek > 6 {
		sc.FirstWeek = 0
	}

	var cds []*CalendarData
	solarDates := sc.SolarCalendar(y, m, d)
	for _, dv := range solarDates {
		cd := new(CalendarData)
		cd.SD = dv
		ld, err := sc.Solar2Lunar(dv)
		if err == nil {
			cd.LD = ld
		}
		cds = append(cds, cd)
	}

	return cds
}

// SolarCalendar 当前时间time.Now() 日历,公历
func (sc *Calendar) NowSolarCalendar() []*SolarDate {
	if sc.Loc == nil {
		sc.Loc = time.Local
	}
	if sc.FirstWeek < 0 || sc.FirstWeek > 6 {
		sc.FirstWeek = 0
	}
	sd := TimeToDate(timeFun().In(sc.Loc))
	return sc.solarCalendar(sd)
}

// SolarCalendar 日历,公历
func (sc *Calendar) SolarCalendar(y, m, d int) []*SolarDate {
	if sc.Loc == nil {
		sc.Loc = time.Local
	}
	if sc.FirstWeek < 0 || sc.FirstWeek > 6 {
		sc.FirstWeek = 0
	}

	hour, min, sec := timeFun().In(sc.Loc).Clock()
	sd := TimeToDate(time.Date(y, time.Month(m), d, hour, min, sec, 0, sc.Loc))

	return sc.solarCalendar(sd)
}

// calendars 公历和农历
func (sc *Calendar) calendars(sd *Date) []*CalendarData {
	var cds []*CalendarData
	solarDates := sc.solarCalendar(sd)
	for _, dv := range solarDates {
		cd := new(CalendarData)
		cd.SD = dv
		ld, err := sc.Solar2Lunar(dv)
		if err == nil {
			cd.LD = ld
		}
		cds = append(cds, cd)
	}

	return cds
}

// solarCalendar 公历日历
func (sc *Calendar) solarCalendar(sd *Date) []*SolarDate {
	var jqmap map[string]*SolarJQ
	var err error
	if sc.Getjq {
		// 取节气
		_, jqmap, err = sc.Jieqi(sd.Year)
		if err != nil {
			jqmap = nil
		}
	}

	if sc.Grid == GridDay {
		insd := new(SolarDate)
		insd.Date = sd.clone()
		sc.tianGanDiZhi(insd)
		if sc.Getjq {
			jqindex := fmt.Sprintf("%d-%d-%d", insd.Date.Year, insd.Date.Month, insd.Date.Day)
			if jqv, ok := jqmap[jqindex]; ok {
				insd.Jq = jqv
			}
		}
		return []*SolarDate{insd}
	} else if sc.Grid == GridWeek {
		return sc.weekCalendar(sd, jqmap)
	} else if sc.Grid == GridMonth {
		var resultDate []*SolarDate
		mFirstDayDate := sd.monthFirstDayDate()
		for i := 0; i < 6; i++ {
			addDays := i * 7
			resultDate = append(resultDate, sc.weekCalendar(mFirstDayDate.addDays(addDays), jqmap)...)
		}

		return resultDate
	}

	return nil
}

// weekCalendar 按周取日历
func (sc *Calendar) weekCalendar(sd *Date, jqmap map[string]*SolarJQ) []*SolarDate {
	var subdays int
	if sd.Week >= sc.FirstWeek {
		subdays = sd.Week - sc.FirstWeek
	} else {
		subdays = 7 - sc.FirstWeek + sd.Week
	}
	firstDate := sd.addDays(-subdays)
	var resultDate []*SolarDate
	for i := 0; i < 7; i++ {
		insd := new(SolarDate)
		insd.Date = firstDate.addDays(i)
		sc.tianGanDiZhi(insd)
		if sc.Getjq {
			jqindex := fmt.Sprintf("%d-%d-%d", insd.Date.Year, insd.Date.Month, insd.Date.Day)
			if jqv, ok := jqmap[jqindex]; ok {
				insd.Jq = jqv
			}
		}
		resultDate = append(resultDate, insd)
	}
	return resultDate
}

// DateToSolarDate 用Date初始SolarDate
func (sc *Calendar) DateToSolarDate(d *Date) *SolarDate {
	sd := new(SolarDate)
	sd.Date = d.clone()
	sc.tianGanDiZhi(sd)
	return sd
}

// Solar2Julian 将公历时间转换为儒略日历时间
func (sc *Calendar) Solar2Julian(sd *Date) (float64, error) {

	if sc.Loc == nil {
		sc.Loc = time.Local
	}

	yy := float64(sd.Year)
	mm := float64(sd.Month)
	dd := float64(sd.Day)
	hh := float64(sd.Hour)
	mi := float64(sd.Min)
	ss := float64(sd.Sec)

	yp := yy + math.Floor((mm-3)/float64(10))

	var init float64 = 0
	var jdy float64

	if (yy > 1582) || (yy == 1582 && mm > 10) || (yy == 1582 && mm == 10 && dd >= 15) {
		init = 1721119.5
		jdy = math.Floor(yp*365.25) - math.Floor(yp/100) + math.Floor(yp/400)
	}
	if (yy < 1582) || (yy == 1582 && mm < 10) || (yy == 1582 && mm == 10 && dd <= 4) {
		init = 1721117.5
		jdy = math.Floor(yp * 365.25)
	}

	// 因历法转换,现用公历中1582年10月5日-14日是不存在的
	if init == 0 {
		return 0, errors.New("公历1582年10月5日-14日被跳过")
	}

	mp := float64(int(mm+9) % 12)
	jdm := mp*30 + math.Floor((mp+1)*34/57)
	jdd := dd - 1
	jdh := (hh + (mi+(ss/60))/60) / 24

	return Round(jdy+jdm+float64(jdd)+float64(jdh)+init, 7), nil
}

// Julian2Solar 将儒略日历时间转换为公历(格里高利历)时间
func (sc *Calendar) Julian2Solar(jd float64) *Date {
	if sc.Loc == nil {
		sc.Loc = time.Local
	}

	var y4h float64
	var init float64

	if jd >= 2299160.5 { // 1582年10月15日,此日起是儒略日历,之前是儒略历
		y4h = 146097
		init = 1721119.5
	} else {
		y4h = 146100
		init = 1721117.5
	}

	jdr := math.Floor(jd - init)

	yh := y4h / 4
	cen := math.Floor((jdr + 0.75) / yh)
	da := math.Floor(jdr + 0.75 - cen*yh)

	ywl := float64(1461) / 4
	jy := math.Floor((da + 0.75) / ywl)
	da = math.Floor(da + 0.75 - ywl*jy + 1)

	ml := float64(153) / 5
	mp := math.Floor((da - 0.5) / ml)
	da = math.Floor((da - 0.5) - 30.6*mp + 1)

	y := (100 * cen) + jy
	m := float64((int(mp)+2)%12) + 1
	if m < 3 {
		y = y + 1
	}
	sd := math.Floor(Round((jd+0.5-math.Floor(jd+0.5))*24, 5)*60*60 + 0.00005)

	mt := math.Floor(sd / 60)
	ss := int(sd) % 60
	hh := math.Floor(mt / 60)
	mt = float64(int(mt) % 60)
	yy := math.Floor(y)
	mm := math.Floor(m)
	dd := math.Floor(da)

	t := time.Date(int(yy), time.Month(int(mm)), int(dd), int(hh), int(mt), int(ss), 0, sc.Loc)
	return TimeToDate(t)
}

// Solar2Lunar 将公历日期转换成农历日期
func (sc *Calendar) Solar2Lunar(sd *SolarDate) (*LunarDate, error) {

	if sc.Loc == nil {
		sc.Loc = time.Local
	}

	// 求出指定年月日之JD值
	s := sd.clone()
	s.Hour = 12
	s.Min = 0
	s.Sec = 0
	jd, err := sc.Solar2Julian(s)
	if err != nil {
		return nil, err
	}
	prev := 0
	_, jdnm, mc, err := sc.zqAndSMandLunarMonthCode(s.Year)

	if err != nil {
		return nil, err
	}

	if math.Floor(jd) < math.Floor(jdnm[0]+0.5) {

		prev = 1
		_, jdnm, mc, err = sc.zqAndSMandLunarMonthCode(s.prevYear().Year)
		if err != nil {
			return nil, err
		}
	}

	var mi = 0
	for i := 0; i <= 14; i++ { // 指令中加0.5是为了改为从0时算起而不从正午算起
		if math.Floor(jd) >= math.Floor(jdnm[i]+0.5) && math.Floor(jd) < math.Floor(jdnm[i+1]+0.5) {
			mi = i
			break
		}
	}

	ld := new(LunarDate)
	ld.Date = sd.clone()

	if mc[mi] < 2 || prev == 1 { // 年
		ld.Year = sd.prevYear().Year
	}

	ld.LeapYear = sc.leap(ld.Year) // 闰几月，0为无闰月

	isLeapMonth := 0 // 初始,该月是否为闰月
	if ld.LeapYear > 0 && (mc[mi]-math.Floor(mc[mi]))*2+1 != 1 { // 因mc(mi)=0对应到前一年农历11月,mc(mi)=1对应到前一年农历12月,mc(mi)=2对应到本年1月,依此类推
		isLeapMonth = 1
	}

	ld.Month = int(math.Floor(mc[mi]+10))%12 + 1 // 月
	ld.MonthStr = MonthChinese(ld.Month)

	if isLeapMonth == 1 {
		ld.LeapMonth = 1 // 当前月是闰月
		ld.LeapStr = LeapStr
	}

	ld.MonthDays, err = sc.lunarDays(ld.Year, ld.Month, isLeapMonth)
	if err != nil {
		return nil, err
	}

	ld.Day = int(math.Floor(jd) - math.Floor(jdnm[mi]+0.5) + 1) // 日,此处加1是因为每月初一从1开始而非从0开始
	ld.DayStr = DayChinese(ld.Day)

	sc.lunarYearGanZiCommon(ld)
	return ld, err
}

// Lunar2Solar 将农历日期转换成公历日期
func (sc *Calendar) Lunar2Solar(ly, lm, ld, isLeap int) (*SolarDate, error) {
	if sc.Loc == nil {
		sc.Loc = time.Local
	}

	if ly < -7000 || ly > 7000 {
		return nil, errors.New("年份超出限制")
	}
	if ly < -1000 || ly > 3000 { // 适用于西元-1000年至西元3000年,超出此范围误差较大
		return nil, errors.New("年份限制西元-1000年至西元3000年")
	}

	if lm < 1 {
		lm = 1
	} else if lm > 12 {
		lm = 12
	}
	if ld < 1 {
		ld = 1
	} else if ld > 29 {
		ldt, err := sc.lunarDays(ly, lm, 0)
		if err != nil {
			return nil, err
		}
		ld = ldt
	}

	_, jdnm, mc, err := sc.zqAndSMandLunarMonthCode(ly)
	if err != nil {
		return nil, err
	}

	leap := 0 // 若闰月旗标为0代表无闰月
	for j := 1; j <= 14; j++ { // 确认指定年前一年11月开始各月是否闰月
		if mc[j]-math.Floor(mc[j]) > 0 { // 若是,则将此闰月代码放入闰月旗标內
			leap = int(math.Floor(mc[j] + 0.5)) // leap=0对应农历11月,1对应农历12月,2对应农历隔年1月,依此类推.
			break
		}
	}

	// 11月对应到1,12月对应到2,1月对应到3,2月对应到4,依此类推
	mm := lm + 2 // 用农历月份简单推出公历月份

	// 求算农历各月之大小,大月30天,小月29天
	nofd := [15]int{}
	for i := 0; i <= 14; i++ {
		nofd[i] = int(math.Floor(jdnm[i+1]+0.5) - math.Floor(jdnm[i]+0.5)) // 每月天数,加0.5是因JD以正午起算
	}

	var jd float64 = 0

	if I2b(isLeap) { // 若是闰月
		if leap < 3 { // 而旗标非闰月或非本年闰月,则表示此年不含闰月.leap=0代表无闰月,=1代表闰月为前一年的11月,=2代表闰月为前一年的12月
			return nil, errors.New("此年非闰年") // 此年非闰年
		} else { // 若本年內有闰月
			if leap != mm { // 但不为指入的月份
				return nil, errors.New("该月非闰月") // 则指定的月份非闰月,此月非闰月
			} else { // 若输入的月份即为闰月
				if lm <= nofd[mm] { // 若指定的日期不大于当月的天數
					jd = jdnm[mm] + float64(ld) - 1 // 则将当月之前的JD值加上日期之前的天數
				} else { // 日期超出范围
					return nil, errors.New("日期超出范围")
				}
			}
		}
	} else {
		if leap == 0 { // 若旗標非閏月,則表示此年不含閏月(包括前一年的11月起之月份)
			if ld <= nofd[mm-1] { // 若輸入的日期不大於當月的天數
				jd = jdnm[mm-1] + float64(ld) - 1 // 則將當月之前的JD值加上日期之前的天數
			} else { // 日期超出范围
				return nil, errors.New("日期超出范围")
			}
		} else { // 若旗標為本年有閏月(包括前一年的11月起之月份) 公式nofd(mx - (mx > leap) - 1)的用意為:若指定月大於閏月,則索引用mx,否則索引用mx-1
			if ld <= nofd[mm+B2i(mm > leap)-1] { // 若輸入的日期不大於當月的天數
				jd = jdnm[mm+B2i(mm > leap)-1] + float64(ld) - 1 // 則將當月之前的JD值加上日期之前的天數
			} else { // 日期超出范围
				return nil, errors.New("日期超出范围")
			}
		}
	}

	d := sc.Julian2Solar(jd)
	ch, cmi, cs := timeFun().In(sc.Loc).Clock()
	d.Hour = ch
	d.Min = cmi
	d.Sec = cs
	sd := new(SolarDate)
	sd.Date = d

	sc.tianGanDiZhi(sd)

	// 取节气
	if sc.Getjq {
		_, jqmap, err := sc.Jieqi(sd.Year)
		if err == nil {
			jqindex := fmt.Sprintf("%d-%d-%d", d.Year, d.Month, d.Day)
			if jqv, ok := jqmap[jqindex]; ok {
				sd.Jq = jqv
			}
		}
	}

	return sd, nil
}

// Jieqi 求出含某公历年立春点开始的24节气
// 这里为了日历好显示，取27个节气，从上一年的冬至开始
// 返回[]*SolarJQ是一个有序切片,map[string]*SolarJQ是一个以"年-月-日"为索引的map
func (sc *Calendar) Jieqi(year int) ([]*SolarJQ, map[string]*SolarJQ, error) {

	if sc.Loc == nil {
		sc.Loc = time.Local
	}

	// jq := []*SolarJQ
	// 为了在日历上显示，加map返回
	var jq []*SolarJQ
	var jqmap = make(map[string]*SolarJQ)

	dj, err := adjustedJQ(year-1, 18, 23)
	if err != nil {
		return nil, nil, err
	}

	ji := -1
	jqak := 0

	for k, v := range dj {
		if k < 18 {
			continue
		}
		if k > 23 {
			continue
		}

		if v == 0 {
			continue
		}

		ji++
		jqak = (ji + 18) % 24 // 节气名称的索引
		jqdate := sc.Julian2Solar(v)
		sjq := &SolarJQ{
			Name: JieQiArray[jqak],
			Date: jqdate,
		}
		jq = append(jq, sjq)

		mindex := fmt.Sprintf("%d-%d-%d", jqdate.Year, jqdate.Month, jqdate.Day)
		jqmap[mindex] = sjq
	}

	dj, err = adjustedJQ(year, 0, 20)
	if err != nil {
		return nil, nil, err
	}

	for _, v := range dj {

		if v == 0 {
			continue
		}

		ji++
		jqak = (ji + 18) % 24
		jqdate := sc.Julian2Solar(v)
		sjq := &SolarJQ{
			Name: JieQiArray[jqak],
			Date: jqdate,
		}
		jq = append(jq, sjq)

		mindex := fmt.Sprintf("%d-%d-%d", jqdate.Year, jqdate.Month, jqdate.Day)
		jqmap[mindex] = sjq
	}

	return jq, jqmap, nil
}

// tianGanDiZhi 公历对应的干支 以立春时间开始
func (sc *Calendar) tianGanDiZhi(sd *SolarDate) {
	sdgz := new(SolarTianGanDiZhi)

	s := sd.clone()

	jd, err := sc.Solar2Julian(s)
	if err != nil {
		return
	}
	jq, err := pureJQsinceSpring(s.Year)
	if err != nil {
		return
	}
	if jd < jq[1] {
		s = s.prevYear()
		jq, err = pureJQsinceSpring(s.Year)
		if err != nil {
			return
		}
	}

	tg := [4]int{}
	dz := [4]int{}
	ygz := ((s.Year+4712+24)%60 + 60) % 60
	sdgz.Ytg = ygz % 10 // 年干
	sdgz.Ydz = ygz % 12 // 年支
	sdgz.YtgStr = TianGanArray[sdgz.Ytg]
	sdgz.YdzStr = DiZhiArray[sdgz.Ydz]

	var ix int
	// 比较求算节气月,求出月干支
	for j := 0; j <= 15; j++ {
		// 已超过指定时刻,故应取前一个节气
		if jq[j] >= jd {
			ix = j - 1
			break
		}
	}

	tmm := ((s.Year+4712)*12 + (ix - 1) + 60) % 60
	mgz := (tmm + 50) % 60
	sdgz.Mtg = mgz % 10 // 月干
	sdgz.Mdz = mgz % 12 // 月支
	sdgz.MtgStr = TianGanArray[sdgz.Mtg]
	sdgz.MdzStr = DiZhiArray[sdgz.Mdz]

	jda := jd + 0.5                                           // 计算日柱之干支,加0.5是将起始点从正午改为从0点开始
	thes := ((jda - math.Floor(jda)) * 86400) + float64(3600) // 将jd的小数部份化为秒,並加上起始点前移的一小时(3600秒),取其整数值
	dayjd := math.Floor(jda) + thes/86400
	dgz := (int(math.Floor(dayjd+49))%60 + 60) % 60
	sdgz.Dtg = dgz % 10 // 日干
	sdgz.Ddz = dgz % 12 // 日支
	// 区分早晚子时,日柱前移一柱
	if sc.Zwz && (s.Hour >= 23) {
		sdgz.Dtg = (tg[2] + 10 - 1) % 10
		sdgz.Ddz = (dz[2] + 12 - 1) % 12
	}
	sdgz.DtgStr = TianGanArray[sdgz.Dtg]
	sdgz.DdzStr = DiZhiArray[sdgz.Ddz]

	dh := dayjd * 12
	hgz := (int(math.Floor(dh+48))%60 + 60) % 60
	sdgz.Htg = hgz % 10
	sdgz.Hdz = hgz % 12
	sdgz.HtgStr = TianGanArray[sdgz.Htg]
	sdgz.HdzStr = DiZhiArray[sdgz.Hdz]

	sd.GanZhi = sdgz
}

// smSinceWinterSolstice 求算以含冬至中气为阴历11月开始的连续16个朔望月
func (sc *Calendar) smSinceWinterSolstice(year int, jdws float64) ([16]float64, error) {

	tjd := [20]float64{}
	jdnm := [16]float64{}

	// 求年初前两个月附近的新月点(即前一年的11月初)
	jd, err := sc.Solar2Julian(YmdNewDate(year, 11, 1, sc.Loc).prevYear())
	if err != nil {
		return jdnm, err
	}

	// 求得自2000年1月起第kn个平均朔望日及其JD值
	// kn,thejd := meanNewMoon(jd)
	kn, _ := meanNewMoon(jd)

	// 求出连续20个朔望月
	for i := 0; i <= 19; i++ {
		k := kn + float64(i)

		// mjd := thejd + synMonth * float64(i)

		// 以k值代入求瞬时朔望日,因中国比格林威治先行8小时,加1/3天
		tjd[i] = trueNewMoon(k) + float64(1)/3

		// 下式为修正 dynamical time to Universal time
		month := i - 1
		tjd[i] = Round(tjd[i]-deltaT(year, month)/1440, 7) // 1为1月,0为前一年12月,-1为前一年11月(当i=0时,i-1=-1,代表前一年11月)
	}

	var kj = 0
	for j := 0; j <= 18; j++ {
		kj = j
		if math.Floor(tjd[j]+0.5) > math.Floor(jdws+0.5) {
			break
		} // 已超过冬至中气(比较日期法)
	}

	/*if kj == 0 {
		kj = 1
	}*/
	for k := 0; k <= 15; k++ { // 取上一步的索引值
		jdnm[k] = tjd[kj-1+k] // 重排索引,使含冬至朔望月的索引为0
	}

	return jdnm, nil

}

// zqAndSMandLunarMonthCode 以比较日期法求算冬月及其余各月名称代码,包含闰月,冬月为0,腊月为1,正月为2,其余类推.闰月多加0.5
func (sc *Calendar) zqAndSMandLunarMonthCode(year int) ([15]float64, [16]float64, [15]float64, error) {

	mc := [15]float64{}
	// 求出自冬至点为起点的连续15个中气
	jdzq, err := zqSinceWinterSolstice(year)
	if err != nil {
		fmt.Println(err.Error())
		return [15]float64{}, [16]float64{}, [15]float64{}, err
	}

	// 求出以含冬至中气为阴历11月(冬月)开始的连续16个朔望月的新月點
	jdnm, err := sc.smSinceWinterSolstice(year, jdzq[0])
	if err != nil {
		return [15]float64{}, [16]float64{}, [15]float64{}, err
	}
	// 设定旗标,0表示未遇到闰月,1表示已遇到闰月
	yz := 0

	if math.Floor(jdzq[12]+0.5) >= math.Floor(jdnm[13]+0.5) {

		for i := 1; i <= 14; i++ {

			// 至少有一个朔望月不含中气,第一个不含中气的月即为闰月
			// 若阴历腊月起始日大於冬至中气日,且阴历正月起始日小于或等于大寒中气日,则此月为闰月,其余同理
			if (jdnm[i]+0.5) > math.Floor(jdzq[i-1-yz]+0.5) && math.Floor(jdnm[i+1]+0.5) <= math.Floor(jdzq[i-yz]+0.5) {
				mc[i] = float64(i) - 0.5
				yz = 1 // 标示遇到闰月
			} else {
				mc[i] = float64(i - yz) // 遇到闰月开始,每个月号要减1
			}

		}

	} else { // 否则表示两个连续冬至之间只有11个整月,故无闰月

		for i := 0; i <= 12; i++ { // 直接赋予这12个月月代码
			mc[i] = float64(i)
		}
		for i := 13; i <= 14; i++ { // 处理次一置月年的11月与12月,亦有可能含闰月
			// 若次一阴历腊月起始日大于附近的冬至中气日,且阴历正月起始日小于或等于大寒中气日,则此月为腊月,次一正月同理.
			if (jdnm[i]+0.5) > math.Floor(jdzq[i-1-yz]+0.5) && math.Floor(jdnm[i+1]+0.5) <= math.Floor(jdzq[i-yz]+0.5) {
				mc[i] = float64(i) - 0.5
				yz = 1 // 标示遇到闰月
			} else {
				mc[i] = float64(i - yz) // 遇到闰月开始,每个月号要减1
			}
		}

	}

	return jdzq, jdnm, mc, nil
}

// leap 获取农历某年的闰月,0为无闰月
func (sc *Calendar) leap(ly int) int {
	_, _, mc, err := sc.zqAndSMandLunarMonthCode(ly)
	if err != nil {
		return 0
	}

	var leap float64 = 0 // 若闰月旗标为0代表无闰月
	for j := 1; j <= 14; j++ { // 确认指定年前一年11月开始各月是否闰月
		if mc[j]-math.Floor(mc[j]) > 0 { // 若是,则将此闰月代码放入闰月旗标內
			leap = math.Floor(mc[j] + 0.5) // leap = 0对应农历11月,1对应农历12月,2对应农历隔年1月,依此类推.
			break
		}
	}

	return int(math.Max(0, leap-2))
}

// lunarDays 获取农历某个月有多少天
func (sc *Calendar) lunarDays(ly, lm, isLeap int) (int, error) {

	_, jdnm, mc, err := sc.zqAndSMandLunarMonthCode(ly)
	if err != nil {
		return 0, err
	}

	leap := 0 // 若闰月旗标为0代表无闰月
	for j := 1; j <= 14; j++ { // 确认指定年前一年11月开始各月是否闰月
		if mc[j]-math.Floor(mc[j]) > 0 { // 若是,则将此闰月代码放入闰月旗标內
			leap = int(math.Floor(mc[j] + 0.5)) // leap=0对应农历11月,1对应农历12月,2对应农历隔年1月,依此类推.
			break
		}
	}

	// 11月对应到1,12月对应到2,1月对应到3,2月对应到4,依此类推
	lm = lm + 2 // 用农历月份简单推出公历月份

	// 求算农历各月之大小,大月30天,小月29天
	nofd := [15]int{}
	for i := 0; i <= 14; i++ {
		nofd[i] = int(math.Floor(jdnm[i+1]+0.5) - math.Floor(jdnm[i]+0.5)) // 每月天数,加0.5是因JD以正午起算
	}

	dy := 0 // 当月天数

	if I2b(isLeap) { // 若是闰月
		if leap < 3 { // 而旗标非闰月或非本年闰月,则表示此年不含闰月.leap=0代表无闰月,=1代表闰月为前一年的11月,=2代表闰月为前一年的12月
			return 0, errors.New("该年非闰年")
		}
		// 若本年內有闰月
		if leap != lm { // 但不为指定的月份
			return 0, errors.New("该月非该年的闰月")
		} else { // 若指定的月份即为闰月
			dy = nofd[lm]
		}
	} else { // 若没有指明是闰月
		if leap == 0 { // 若旗标非闰月,则表示此年不含闰月(包括前一年的11月起之月份)
			dy = nofd[lm-1]
		} else { // 若旗标为本年有闰月(包括前一年的11月起之月份) 公式nofd(mx - (mx > leap) - 1)的用意为:若指定月大于闰月,则索引用mx,否则索引用mx-1
			dy = nofd[lm+B2i(lm > leap)-1]
		}
	}

	return dy, nil
}

// 农历年干支通俗记法，以春节开始
func (sc *Calendar) lunarYearGanZiCommon(ld *LunarDate) {
	gk := (ld.Year - 3) % 10
	zk := (ld.Year - 3) % 12

	if gk == 0 {
		gk = 9
	} else {
		gk -= 1
	}

	if zk == 0 {
		zk = 11
	} else {
		zk -= 1
	}
	ld.YearGanZi = &LunarYearGanZi{
		Gan:    TianGanArray[gk],
		Zhi:    DiZhiArray[zk],
		Animal: SymbolicAnimalsArray[zk],
	}
}

// WeekChinese 星期中文
func WeekChinese(w int) string {
	if w < 0 || w > 6 {
		return ""
	}
	return NumberChineseArray[w]
}

// MonthChinese 农历月份常用名称
func MonthChinese(m int) string {
	if m > 0 && m <= 12 {
		return MonthChineseArray[m-1]
	}
	return ""
}

// DayChinese 农历日期数字返回汉字表示法
func DayChinese(d int) string {
	daystr := ""
	if d < 1 || d > 30 {
		return ""
	}
	switch d {
	case 10:
		daystr = DayChineseArray[0] + NumberChineseArray[10]
	case 20:
		daystr = DayChineseArray[2] + NumberChineseArray[10]
	case 30:
		daystr = DayChineseArray[3] + NumberChineseArray[10]
	default:
		k := d / 10
		m := d % 10
		daystr = DayChineseArray[k] + NumberChineseArray[m]
	}
	return daystr
}

// perturbation 地球在绕日运行时会因受到其他星球之影响而产生摄动(perturbation)
// 返回某时刻(儒略日历)的摄动偏移量
func perturbation(jd float64) float64 {
	t := (jd - 2451545) / 36525
	var s float64 = 0
	for k := 0; k <= 23; k++ {
		s += ptsa[k] * math.Cos(ptsb[k]*2*math.Pi/360+ptsc[k]*2*math.Pi/360*t)
	}
	w := 35999.373*t - 2.47
	l := 1 + 0.0334*math.Cos(w*2*math.Pi/360) + 0.0007*math.Cos(2*w*2*math.Pi/360)
	return Round(0.00001*s/l, 16)
}

// trueNewMoon 求出实际新月点
// 以2000年初的第一个均值新月点为0点求出的均值新月点和其朔望月之序數 k 代入此副程式來求算实际新月点
func trueNewMoon(k float64) float64 {

	jdt := bnm + k*synMonth

	t := (jdt - 2451545) / 36525 // 2451545为2000年1月1日正午12时的JD

	// t2 := t * t // square for frequent use
	t2 := math.Pow(t, 2)

	// t3 := t2 * t // cube for frequent use
	t3 := math.Pow(t, 3)

	// t4 := t3 * t // to the fourth
	t4 := math.Pow(t, 4)

	// mean time of phase
	// 加上调整值0.0001337*t2 - 0.00000015*t3 + 0.00000000073*t4
	pt := jdt + 0.0001337*t2 - 0.00000015*t3 + 0.00000000073*t4
	// Sun's mean anomaly(地球绕太阳运行均值近点角)(从太阳观察)
	m := 2.5534 + 29.10535669*k - 0.0000218*t2 - 0.00000011*t3
	// Moon's mean anomaly(月球绕地球运行均值近点角)(从地球观察)
	mprime := 201.5643 + 385.81693528*k + 0.0107438*t2 + 0.00001239*t3 - 0.000000058*t4
	// Moon's argument of latitude(月球的纬度参数)
	f := 160.7108 + 390.67050274*k - 0.0016341*t2 - 0.00000227*t3 + 0.000000011*t4
	// Longitude of the ascending node of the lunar orbit(月球绕日运行轨道升交点之经度)
	omega := 124.7746 - 1.5637558*k + 0.0020691*t2 + 0.00000215*t3
	// 乘式因子
	es := 1 - 0.002516*t - 0.0000074*t2
	// 因perturbation造成的偏移
	pi180 := math.Pi / 180
	apt1 := -0.4072 * math.Sin(pi180*mprime)
	apt1 += 0.17241 * es * math.Sin(pi180*m)
	apt1 += 0.01608 * math.Sin(pi180*2*mprime)
	apt1 += 0.01039 * math.Sin(pi180*2*f)
	apt1 += 0.00739 * es * math.Sin(pi180*(mprime-m))
	apt1 -= 0.00514 * es * math.Sin(pi180*(mprime+m))
	apt1 += 0.00208 * es * es * math.Sin(pi180*(2*m))
	apt1 -= 0.00111 * math.Sin(pi180*(mprime-2*f))
	apt1 -= 0.00057 * math.Sin(pi180*(mprime+2*f))
	apt1 += 0.00056 * es * math.Sin(pi180*(2*mprime+m))
	apt1 -= 0.00042 * math.Sin(pi180*3*mprime)
	apt1 += 0.00042 * es * math.Sin(pi180*(m+2*f))
	apt1 += 0.00038 * es * math.Sin(pi180*(m-2*f))
	apt1 -= 0.00024 * es * math.Sin(pi180*(2*mprime-m))
	apt1 -= 0.00017 * math.Sin(pi180*omega)
	apt1 -= 0.00007 * math.Sin(pi180*(mprime+2*m))
	apt1 += 0.00004 * math.Sin(pi180*(2*mprime-2*f))
	apt1 += 0.00004 * math.Sin(pi180*(3*m))
	apt1 += 0.00003 * math.Sin(pi180*(mprime+m-2*f))
	apt1 += 0.00003 * math.Sin(pi180*(2*mprime+2*f))
	apt1 -= 0.00003 * math.Sin(pi180*(mprime+m+2*f))
	apt1 += 0.00003 * math.Sin(pi180*(mprime-m+2*f))
	apt1 -= 0.00002 * math.Sin(pi180*(mprime-m-2*f))
	apt1 -= 0.00002 * math.Sin(pi180*(3*mprime+m))
	apt1 += 0.00002 * math.Sin(pi180*(4*mprime))

	apt2 := 0.000325 * math.Sin(pi180*(299.77+0.107408*k-0.009173*t2))
	apt2 += 0.000165 * math.Sin(pi180*(251.88+0.016321*k))
	apt2 += 0.000164 * math.Sin(pi180*(251.83+26.651886*k))
	apt2 += 0.000126 * math.Sin(pi180*(349.42+36.412478*k))
	apt2 += 0.00011 * math.Sin(pi180*(84.66+18.206239*k))
	apt2 += 0.000062 * math.Sin(pi180*(141.74+53.303771*k))
	apt2 += 0.00006 * math.Sin(pi180*(207.14+2.453732*k))
	apt2 += 0.000056 * math.Sin(pi180*(154.84+7.30686*k))
	apt2 += 0.000047 * math.Sin(pi180*(34.52+27.261239*k))
	apt2 += 0.000042 * math.Sin(pi180*(207.19+0.121824*k))
	apt2 += 0.00004 * math.Sin(pi180*(291.34+1.844379*k))
	apt2 += 0.000037 * math.Sin(pi180*(161.72+24.198154*k))
	apt2 += 0.000035 * math.Sin(pi180*(239.56+25.513099*k))
	apt2 += 0.000023 * math.Sin(pi180*(331.55+3.592518*k))

	return pt + apt1 + apt2
}

// meanNewMoon 对于指定日期时刻所属的朔望月,求出其均值新月点的月序数
func meanNewMoon(jd float64) (float64, float64) {

	// kn为从2000年1月6日14时20分36秒起至指定年月日之农历月数,以synodic month为单位
	// bnm=2451550.09765为2000年1月6日14时20分36秒之JD值.
	kn := math.Floor((jd - bnm) / synMonth)

	jdt := bnm + kn*synMonth

	// Time in Julian centuries from 2000 January 0.5.
	t := (jdt - 2451545) / 36525
	thejd := jdt + 0.0001337*math.Pow(t, 2) - 0.00000015*math.Pow(t, 3) + 0.00000000073*math.Pow(t, 4)

	return kn, thejd
}

// ve 计算指定年(公历)的春分点(mean vernal equinox)
// 比利时的气象学家Jean Meeus在1991年出版的”Astronomical Algorithms”一书中提供了一些求均值春分点(mean vernal equinox)的公式
// 但因地球在绕日运行时会因受到其他星球之影响而产生摄动(perturbation),必须将此现象产生的偏移量加入
// 返回儒略日历格林威治时间
func ve(year int) (float64, error) {

	if year < -8000 {
		return 0, errors.New("年份超出限制")
	}
	if year > 8001 {
		return 0, errors.New("年份超出限制")
	}

	if year >= 1000 && year <= 8001 {
		m := (float64(year) - 2000) / 1000
		s := 2451623.80984 + 365242.37404*m + 0.05169*math.Pow(m, 2) - 0.00411*math.Pow(m, 3) - 0.00057*math.Pow(m, 4)
		return s, nil
	}

	if year >= -8000 && year < 1000 {
		m := float64(year) / 1000
		s := 1721139.29189 + 365242.1374*m + 0.06134*math.Pow(m, 2) + 0.00111*math.Pow(m, 3) - 0.00071*math.Pow(m, 4)
		return s, nil
	}

	return 0, errors.New("年份超出限制")
}

// deltaT 求∆t
func deltaT(year, month int) float64 {

	y := float64(year) + (float64(month)-0.5)/12

	var dt float64

	switch {
	case y <= -500:
		t := (y - 1820) / 100
		dt = -20 + 32*math.Pow(t, 2)
	case y < 500:
		t := y / 100
		dt = 10583.6 - 1014.41*t + 33.78311*math.Pow(t, 2) - 5.952053*math.Pow(t, 3) - 0.1798452*math.Pow(t, 4) + 0.022174192*math.Pow(t, 5) + 0.0090316521*math.Pow(t, 6)
	case y < 1600:
		t := (y - 1000) / 100
		dt = 1574.2 - 556.01*t + 71.23472*math.Pow(t, 2) + 0.319781*math.Pow(t, 3) - 0.8503463*math.Pow(t, 4) - 0.005050998*math.Pow(t, 5) + 0.0083572073*math.Pow(t, 6)
	case y < 1700:
		t := y - 1600
		dt = 120 - 0.9808*t - 0.01532*math.Pow(t, 2) + math.Pow(t, 3)/7129
	case y < 1800:
		t := y - 1700
		dt = 8.83 + 0.1603*t - 0.0059285*math.Pow(t, 2) + 0.00013336*math.Pow(t, 3) - math.Pow(t, 4)/1174000
	case y < 1860:
		t := y - 1800
		dt = 13.72 - 0.332447*t + 0.0068612*math.Pow(t, 2) + 0.0041116*math.Pow(t, 3) - 0.00037436*math.Pow(t, 4) + 0.0000121272*math.Pow(t, 5) - 0.0000001699*math.Pow(t, 6) + 0.000000000875*math.Pow(t, 7)
	case y < 1900:
		t := y - 1860
		dt = 7.62 + 0.5737*t - 0.251754*math.Pow(t, 2) + 0.01680668*math.Pow(t, 3) - 0.0004473624*math.Pow(t, 4) + math.Pow(t, 5)/233174
	case y < 1920:
		t := y - 1900
		dt = -2.79 + 1.494119*t - 0.0598939*math.Pow(t, 2) + 0.0061966*math.Pow(t, 3) - 0.000197*math.Pow(t, 4)
	case y < 1941:
		t := y - 1920
		dt = 21.2 + 0.84493*t - 0.0761*math.Pow(t, 2) + 0.0020936*math.Pow(t, 3)
	case y < 1961:
		t := y - 1950
		dt = 29.07 + 0.407*t - math.Pow(t, 2)/233 + math.Pow(t, 3)/2547
	case y < 1986:
		t := y - 1975
		dt = 45.45 + 1.067*t - math.Pow(t, 2)/260 - math.Pow(t, 3)/718
	case y < 2005:
		t := y - 2000
		dt = 63.86 + 0.3345*t - 0.060374*math.Pow(t, 2) + 0.0017275*math.Pow(t, 3) + 0.000651814*math.Pow(t, 4) + 0.00002373599*math.Pow(t, 5)
	case y < 2050:
		t := y - 2000
		dt = 62.92 + 0.32217*t + 0.005589*math.Pow(t, 2)
	case y < 2150:
		t := (y - 1820) / 100
		dt = -20 + 32*math.Pow(t, 2) - 0.5628*(2150-y)
	default:
		t := (y - 1820) / 100
		dt = -20 + 32*math.Pow(t, 2)
	}

	if y < 1955 || y >= 2005 {
		dt = dt - (0.000012932 * (y - 1955) * (y - 1955))
	}
	return Round(dt/60, 13) // 将秒转换为分
}

// meanJqJd 获取指定年的春分开始的24节气,另外多取2个确保覆盖完一个公历年
// 大致原理是:先用此方法得到理论值,再用摄动值(Perturbation)和固定参数deltaT做调整
func meanJqJd(year int) ([26]float64, error) {

	// 另外多取2个确保覆盖完一个公历年
	// num := 24 + 2

	jqjd := [26]float64{}

	// 该年的春分點

	jd, err := ve(year)
	if err != nil {
		return [26]float64{}, err
	}
	// 该年的回归年长
	ntjd, err := ve(year + 1)
	if err != nil {
		return jqjd, err
	}
	ty := ntjd - jd

	ath := 2 * math.Pi / 24

	tx := (jd - 2451545) / 365250

	e := 0.0167086342 - 0.0004203654*tx - 0.0000126734*tx*tx + 0.0000001444*tx*tx*tx - 0.0000000002*tx*tx*tx*tx + 0.0000000003*tx*tx*tx*tx*tx
	tt := float64(year) / 1000

	vp := 111.25586939 - 17.0119934518333*tt - 0.044091890166673*tt*tt - 4.37356166661345E-04*tt*tt*tt + 8.16716666602386E-06*tt*tt*tt*tt
	rvp := vp * 2 * math.Pi / 360

	var peri [26]float64

	for i := 0; i < cap(peri); i++ {
		flag := 0
		th := ath*float64(i) + rvp

		if th > math.Pi && th <= 3*math.Pi {
			th = 2*math.Pi - th
			flag = 1
		}

		if th > 3*math.Pi {
			th = 4*math.Pi - th
			flag = 2
		}

		f1 := 2 * math.Atan(math.Sqrt((1-e)/(1+e))*math.Tan(th/2))

		f2 := (e * math.Sqrt(1-e*e) * math.Sin(th)) / (1 + e*math.Cos(th))

		f := (f1 - f2) * ty / 2 / math.Pi

		if flag == 1 {
			f = ty - f
		}
		if flag == 2 {
			f = 2*ty - f
		}
		peri[i] = f
	}
	for i := 0; i < cap(peri); i++ {
		jqjd[i] = Round(jd+peri[i]-peri[0], 8)
	}
	return jqjd, nil
}

// adjustedJQ 获取指定年的春分开始作Perturbaton调整后的24节气,可以多取2个
func adjustedJQ(year, start, end int) ([26]float64, error) {

	jq := [26]float64{}

	// 获取该年春分开始的24节气时间点
	jqjd, err := meanJqJd(year)
	if err != nil {
		return [26]float64{}, err
	}

	for k, jd := range jqjd {
		if k < start {
			continue
		}
		if k > end {
			continue
		}

		// 取得受perturbation影响所需微调
		ptb := perturbation(jd)

		// 修正dynamical time to Universal time
		month := int(math.Floor((float64(k)+1)/2) + 3)
		dt := deltaT(year, month)
		// 因中国(北京、重庆、上海)时间比格林威治时间先行8小时，即1/3日
		jq[k] = Round(jd+ptb-dt/60/24+float64(1)/3, 8)
	}

	return jq, nil
}

// pureJQsinceSpring 求出以某年立春点开始的节(注意:为了方便计算起运数,此处第0位为上一年的小寒)
func pureJQsinceSpring(year int) ([16]float64, error) {
	jdpjq := [16]float64{}

	dj, err := adjustedJQ(year-1, 19, 23)
	if err != nil {
		return jdpjq, err
	}

	ki := -1 // 数组索引
	for k, v := range dj {
		if k < 19 {
			continue
		}
		if k > 23 {
			continue
		}
		if k%2 == 0 {
			continue
		}

		ki++
		jdpjq[ki] = v
	}

	dj, err = adjustedJQ(year, 0, 25)
	if err != nil {
		return [16]float64{}, err
	}

	for k, v := range dj {
		if k%2 == 0 || v == 0 {
			continue
		}

		ki++
		jdpjq[ki] = v
	}

	return jdpjq, err
}

// zqSinceWinterSolstice 求出自冬至点为起点的连续15个中气
func zqSinceWinterSolstice(year int) ([15]float64, error) {
	jdpjq := [15]float64{}
	dj, err := adjustedJQ(year-1, 18, 23)
	if err != nil {
		return jdpjq, err
	}

	jdpjq[0] = dj[18] // 冬至
	jdpjq[1] = dj[20] // 大寒
	jdpjq[2] = dj[22] // 雨水

	// 求出指定年节气之JD值
	dj, err = adjustedJQ(year, 0, 23)
	if err != nil {
		return [15]float64{}, err
	}

	ki := 2 // 数组索引
	for k, v := range dj {
		if k%2 != 0 || v == 0 {
			continue
		}

		ki++
		jdpjq[ki] = v
	}
	return jdpjq, err
}
