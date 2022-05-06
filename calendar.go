/**
一个用golang写的日历，有公历转农历，农历转公历，节气，干支，星座，生肖等功能
中国的农历历法综合了太阳历和月亮历,为中国的生活生产提供了重要的帮助,是中国古人智慧与中国传统文化的一个重要体现

程序比较准确的计算出农历与二十四节气(精确到分),时间限制在-1000至3000年间,在实际使用中注意限制年份
*/

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

const Author = "liujiawm@gmail.com"
const Version = "1.1.1.211225_beta"

// type SolarTermItem struct 节气
type SolarTermItem struct {
	Index int        `json:"index"` // 节气索引
	Name  string     `json:"name"`  // 节气名称
	Time  *time.Time `json:"time"`  // 定节气时间
}

// type yearSolarTermTemp struct 节气缓存年表
type yearSolarTermTemp struct {
	data map[int][]*SolarTermItem
	mu   sync.RWMutex
}

// type StarSignItem struct 星座单元
type StarSignItem struct {
	Index int    `json:"index"` // 星座索引
	Name  string `json:"name"`  // 星座名称
}

// type FestivalItem struct 节日单元
type FestivalItem struct {
	Show      []string `json:"show"` // 在日历表上显示的节日
	Secondary []string `json:"scdr"` // 其它次要节日
}

// type yearFestivalTemp struct 节日缓存年表
type yearFestivalTemp struct {
	data map[int]map[string][]string
	mu   sync.RWMutex
}

// type CalendarItem struct 日历单元
type CalendarItem struct {
	Time         *time.Time     `json:"time"`     // 格里高历(公历)时间
	IsAccidental int            `json:"isam"`     // 0为本月日期,-1为上一个月的日期,1为下一个月的日期,
	IsToday      int            `json:"istoday"`  // 是否是今天,0不是,1是
	Festival     *FestivalItem  `json:"festival"` // 公历节日
	SolarTerm    *SolarTermItem `json:"st"`       // 节气
	GZ           *GZ            `json:"gz"`       // 干支
	LunarDate    *LunarDate     `json:"ld"`       // 农历
	StarSign     *StarSignItem  `json:"ss"`       // 星座
}

// Calendar的一些临时数据
type CalendarTempData struct {
	st  *yearSolarTermTemp    // 一整年的节气(24)加上一年最后一个和下一年第一个，共26个节气
	jSS *pureJieQi16Temp      // 对应公历某年立春点开始的16个节
	qSS *pureJieQi16Temp      // 对应公历某年冬至开始的16个中气
	tNM *trueNewMoon20Temp    // 以应公历某年连续20个朔望月
	lMC *lunarMonthCode15Temp // 对应农历某年的月份表
	lMD *lunarMonthDays15Temp // 对应农历某年的月份对应天数表
	lFD *yearFestivalTemp     // 对应农历某年的节日表
	gFD *yearFestivalTemp     // 对应公历某年的节日表
}

// 初始Calendar的临时数据
func newCalendarTempData() *CalendarTempData {
	return &CalendarTempData{
		st:  new(yearSolarTermTemp),
		jSS: new(pureJieQi16Temp),
		qSS: new(pureJieQi16Temp),
		tNM: new(trueNewMoon20Temp),
		lMC: new(lunarMonthCode15Temp),
		lMD: new(lunarMonthDays15Temp),
		lFD: new(yearFestivalTemp),
		gFD: new(yearFestivalTemp),
	}
}

// type Calendar struct 日历表
type Calendar struct {
	Items    []*CalendarItem   `json:"items"` // 日历表中所有的日期单元
	config   *CalendarConfig                  // 配置
	loc      *time.Location                   // time.Location 默认time.Local
	rawTime  *time.Time                       // 初始时间,指定的时间
	tempData *CalendarTempData                // 缓存数据
}

var (
	lunarLeapString = "闰" // 农历闰月标志

	// 星期
	weekNameArray = [7]string{"日", "一", "二", "三", "四", "五", "六"}

	// 农历整十日期常用称呼
	lunarWholeTensArray = [4]string{"初", "十", "廿", "卅"}

	// 农历日期数字
	lunarNumberArray = [11]string{"日", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}

	// 农历月份常用称呼
	lunarMonthNameArray = [12]string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "腊"}

	// 天干
	heavenlyStemsNameArray = [10]string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

	// 地支
	earthlyBranchesNameArray = [12]string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 生肖
	symbolicAnimalsNameArray = [12]string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}

	// 星座名称
	starSignsNameArray = [12]string{"水瓶", "双鱼", "白羊", "金牛", "双子", "巨蟹", "狮子", "处女", "天秤", "天蝎", "射手", "摩羯"}

	// 节气名称
	solarTermsNameArray = [24]string{"春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露",
		"秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至", "小寒", "大寒", "立春", "雨水", "惊蛰"}

	// 农历节日,map的索引M表示月,D表示日,如8M15D表示8月15日,节日名称前加"*"号表示重要且在日历表上显示,同一天多个节日用","分隔
	// 支持两个特殊字符$和@, 用M$表示该月最后一日,如12M$ 表示12月最后一日， 用@M表示闰月,如 5@M12D 表示闰5月12日.
	lunarFestivalArray = map[string]string{"1M1D": "*春节", "1M15D": "*元宵节", "5M5D": "*端午节", "7M7D": "*七夕节", "7M15D": "*中元节",
		"8M15D": "*中秋节", "9M9D": "*重阳节", "12M8D": "*腊八节", "12M24D": "*小年", "12M$": "*除夕"}

	// 公历节日,map的索引M表示月,D表示日,如2M14D表示2月14日,节日名称前加"*"号表示重要且在日历表上显示,同一天多个节日用","分隔
	// 支持某月第几个周几这种方式定义节日,如"5M2W0"表示5月第2个周日;0周日,1周一,2周二...6周六
	gregorianFestivalArray = map[string]string{"1M1D": "*元旦", "2M14D": "*情人节", "3M8D": "*妇女节", "3M12D": "*植树节", "3M15D": "世界消费者权益日",
		"4M1D": "*愚人节", "5M1D": "*劳动节", "5M4D": "*青年节", "5M12D": "*护士节", "5M2W0": "*母亲节", "5M31D": "世界无烟日", "6M1D": "*儿童节", "6M5D": "世界环境日",
		"6M3W0": "*父亲节", "6M26D": "国际禁毒日", "7M7D": "抗战纪念日", "9M10D": "*教师节", "10M1D": "*国庆节", "11M1D": "*万圣节", "12M1D": "世界爱滋病日",
		"12M25D": "*圣诞节"}
)

// DefaultCalendar 默认日历设置
func DefaultCalendar() *Calendar {

	return NewCalendar(defaultConfig())
}

// NewCalendar 日历设置
func NewCalendar(cfg CalendarConfig) *Calendar {

	cfg.Grid = int(math.Mod(math.Abs(float64(cfg.Grid)), 3))
	cfg.FirstWeek = int(math.Mod(math.Abs(float64(cfg.FirstWeek)), 7))

	// 默认时区
	var loc *time.Location

	// TimeZoneName的有效性
	if cfg.TimeZoneName == "" {
		loc = time.Local
		cfg.TimeZoneName = loc.String()
	} else {
		var err error
		loc, err = time.LoadLocation(cfg.TimeZoneName)
		if err != nil {
			loc = time.Local
			cfg.TimeZoneName = loc.String()
		}
	}

	// 默认 rawTime
	rawTime := time.Now().In(loc)

	return &Calendar{
		Items:    nil,
		config:   &cfg,
		loc:      loc,
		rawTime:  &rawTime,
		tempData: newCalendarTempData(),
	}
}

// (*Calendar) SetRawTime 设置rawTime
//
// 该方法返回的*Calendar是清除与rawTime相关值的c *Calendar,这样做是为了支持链接使用,
// 因此该方法必需在相关取值之前使用
func (c *Calendar) SetRawTime(year, month, day int, timeParts ...int) *Calendar {
	var hour, minute,second, millisecond = 0, 0, 0, 0
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

	t := time.Date(year, time.Month(month), day, hour, minute, second, millisecond, c.loc)

	return c.setRawTime(t)
}

// (*Calendar) setRawTime 设置rawTime
func (c *Calendar) setRawTime(t time.Time) *Calendar {
	rawYear := c.GetRawTime().Year()

	t = t.In(c.loc) // 使用c.loc重新设置t的时区
	c.rawTime = &t

	// 重新设置rawTime后，请除原相关数据
	c.Items = nil

	// 清临时数据
	if rawYear != t.Year() {
		c.tempData = newCalendarTempData()
	}

	return c
}

// (*Calendar) GetRawTime 取出rawTime
//
// 返回的将是c.rawTime的一个clone
func (c *Calendar) GetRawTime() time.Time {
	if c.rawTime == nil {
		now := time.Now()
		c.rawTime = &now
	}

	return c.rawTime.AddDate(0, 0, 0) // 使用Time.AddDate(0,0,0)的目的是为了clone一个时间
}

// (*Calendar) GenerateWithDate 用指定的年月日时分秒生成日历表
func (c *Calendar) GenerateWithDate(year, month, day int, timeParts ...int) []*CalendarItem {
	return c.SetRawTime(year, month, day, timeParts...).Generate()
}

// (*Calendar) NextMonth 下一个月
func (c *Calendar) NextMonth() []*CalendarItem {
	t := c.GetRawTime().AddDate(0,1,0)
	c.setRawTime(t)
	return c.Generate()
}

// (*Calendar) PreviousMonth 上一个月
func (c *Calendar) PreviousMonth() []*CalendarItem {
	t := c.GetRawTime().AddDate(0,-1,0)
	c.setRawTime(t)
	return c.Generate()
}

// (*Calendar) NextYear 下一年
func (c *Calendar) NextYear() []*CalendarItem {
	t := c.GetRawTime().AddDate(1,0,0)
	c.setRawTime(t)
	return c.Generate()
}

// (*Calendar) PreviousYear 上一年
func (c *Calendar) PreviousYear() []*CalendarItem {
	t := c.GetRawTime().AddDate(-1,0,0)
	c.setRawTime(t)
	return c.Generate()
}

// (*Calendar) Generate 生成日历表
func (c *Calendar) Generate() []*CalendarItem {
	grid := c.config.Grid
	var result []*CalendarItem
	switch grid {
	case GridDay:
		// 天日历表
		rt := c.GetRawTime()
		ry := rt.Year()
		rm := int(rt.Month())
		ci := c.createItem(rt,ry,rm)
		result = append(result,ci)
	case GridWeek:
		// 周日历表
		result = c.weekCalendar()
	case GridMonth:
		// 月日历表
		result = c.monthCalendar()
	default:
		// 月日历表
		result = c.monthCalendar()
	}

	return result
}

// (*Calendar) weekGregorianCalendar 一周的日历表
func (c *Calendar) weekCalendar() []*CalendarItem {

	rawTime := c.GetRawTime()

	// 本月
	currentMonth := int(rawTime.Month())
	currentYear := rawTime.Year()

	// 日历表首日
	_, firstDayTime := c.firstDay(rawTime)

	// item
	var itemsArray [7]*CalendarItem
	var wg = sync.WaitGroup{}
	wg.Add(cap(itemsArray))
	for i := 0; i < cap(itemsArray); i++ {
		go func(i int) {
			defer wg.Done()

			gt := firstDayTime.AddDate(0, 0, i)
			itemsArray[i] = c.createItem(gt, currentYear, currentMonth)
		}(i)
	}
	wg.Wait()

	var itemsSlice []*CalendarItem
	itemsSlice = itemsArray[:]

	// 附值给Calendar.Items
	c.Items = itemsSlice

	// copy一个返回
	var resultItems = make([]*CalendarItem, len(itemsSlice))
	copy(resultItems, itemsSlice)

	return itemsSlice
}

// (*Calendar) monthGregorianCalendar 一个月的日历表
func (c *Calendar) monthCalendar() []*CalendarItem {

	rawTime := c.GetRawTime()

	// 本月第一天time
	t := BeginningOfMonth(rawTime)

	// 本月
	currentMonth := int(t.Month())
	currentYear := t.Year()

	// 日历表首日
	_, firstDayTime := c.firstDay(t)

	// item
	var itemsArray [42]*CalendarItem
	var wg = sync.WaitGroup{}
	wg.Add(cap(itemsArray))
	for i := 0; i < cap(itemsArray); i++ {
		go func(i int) {
			defer wg.Done()

			gt := firstDayTime.AddDate(0, 0, i)
			itemsArray[i] = c.createItem(gt, currentYear, currentMonth)
		}(i)
	}
	wg.Wait()

	var itemsSlice []*CalendarItem
	itemsSlice = itemsArray[:]

	// 附值给Calendar.Items
	c.Items = itemsSlice

	// copy一个返回
	var resultItems = make([]*CalendarItem, len(itemsSlice))
	copy(resultItems, itemsSlice)

	return itemsSlice
}

// (*Calendar) firstDate 计算t与c.config.FirstWeek相关几日，并同时返回首日time
func (c *Calendar) firstDay(t time.Time) (int, time.Time) {
	// t与日历表首日相差几日，该值根据c.config.FirstWeek计算得出
	var differenceDays int

	// 该日历表的首日
	var firstDayTime time.Time

	// t是周几
	tDayWeek := int(t.Weekday())

	// 根据c.config.FirstWeek算出该日历表的首日
	if tDayWeek >= c.config.FirstWeek {
		differenceDays = tDayWeek - c.config.FirstWeek
	} else {
		differenceDays = 7 - c.config.FirstWeek + tDayWeek
	}

	firstDayTime = t.AddDate(0, 0, -differenceDays)

	return differenceDays, firstDayTime
}

// (*Calendar) createItem 用t计算出单元其它相关值
func (c *Calendar) createItem(t time.Time, currentYear, currentMonth int) *CalendarItem {

	year, _month, day := t.Date()
	month := int(_month)

	item := new(CalendarItem)

	item.Time = &t

	var wg = sync.WaitGroup{}
	wg.Add(7) // 在修改时要注意这里定义goroutine次数

	// 是否非本月的日期,0是本月日期,-1为上一月日期,1为下一月日期
	go func() {
		defer wg.Done()

		isAccidental := 0
		if currentYear > year {
			isAccidental = -1
		} else if currentYear < year {
			isAccidental = 1
		} else {
			if currentMonth > month {
				isAccidental = -1
			} else if currentMonth < month {
				isAccidental = 1
			}
		}
		item.IsAccidental = isAccidental
	}()

	// 是否是今天
	go func() {
		defer wg.Done()

		now := time.Now().In(c.loc)
		nY, nM, nD := now.Date()

		if nD != day || nM != _month || nY != year {
			item.IsToday = 0
		} else {
			item.IsToday = 1
		}

	}()

	// 公历节日
	go func() {
		defer wg.Done()

		gf := c.gregorianFestival(t)
		item.Festival = &gf
	}()

	// 节气
	go func() {
		defer wg.Done()

		if c.config.SolarTerms {
			// sts := c.SolarTerms(c.rawTime.Year())
			sts := c.SolarTerms(year)

			stkStr := "2006-1-2"
			for _, stv := range sts {
				if t.Format(stkStr) == stv.Time.Format(stkStr) {
					item.SolarTerm = stv
					break
				}
			}
		}
	}()

	// 干支
	go func() {
		defer wg.Done()

		if c.config.HeavenlyEarthly {
			gz := c.ChineseSexagenaryCycle(t)
			item.GZ = &gz
		}
	}()

	// 农历
	go func() {
		defer wg.Done()

		if c.config.Lunar {
			ld := c.gregorianToLunar(t, true)
			item.LunarDate = &ld
		}
	}()

	// 星座
	go func() {
		defer wg.Done()

		if c.config.StarSign {
			ssi, ss, _ := StarSign(month, day)
			item.StarSign = &StarSignItem{Index: ssi, Name: ss}
		}
	}()

	wg.Wait()

	return item
}

// (*Calendar) gregorianFestival 取公历节日
func (c *Calendar) gregorianFestival(t time.Time) FestivalItem {
	gregorianYear, gregorianMonth, gregorianDay := t.Date()

	fds := c.tempData.gFD.getData(gregorianYear)

	if len(fds) == 0 {

		for gfK, gfV := range gregorianFestivalArray {
			var err error
			trueK := "" // 经过处理转换成月日的K
			m := ""     // 月
			d := ""     // 日
			n := ""     // 第几周
			w := ""     // 周几
			month := 0
			day := 0

			// 某月第几周几 索引正则
			if strings.Index(gfK, "W") > -1 {
				re := regexp.MustCompile("^([0-9]{1,2})M(?:D)?([1-4])W([0-6])$").FindStringSubmatch(gfK)
				if len(re) != 4 {
					continue // 索引格式不正确
				}
				// 月数
				if re[1] == "" {
					continue // 索引格式不正确,没指明月份
				}
				m = re[1]
				// 月数string转int
				month, err = strconv.Atoi(m)
				if err != nil || month < 1 || month > 12 {
					continue // 月份为空或数字不正确
				}

				// 第几周
				if re[2] == "" {
					continue // 索引格式不正确,没指明第几周
				}
				n = re[2]
				// 第几周string转int
				num, err := strconv.Atoi(n)
				if err != nil || num < 1 || num > 4 {
					continue // 第几周为空或数字不正确
				}

				// 周几
				if re[3] == "" {
					continue // 索引格式不正确,没指明周几
				}
				w = re[3]
				// 周几string转int
				week, err := strconv.Atoi(w)
				if err != nil || week < 0 || week > 6 {
					continue // 第几周为空或数字不正确
				}

				tw := time.Date(gregorianYear, time.Month(month), 1, 0, 0, 0, 100, c.loc)
				tww := int(tw.Weekday())

				differenceDays := 0
				if week >= tww {
					differenceDays = week - tww
				} else {
					differenceDays = 7 - (tww - week)
				}

				day = (num-1)*7 + 1 + differenceDays

				if GregorianMonthDays(gregorianYear, month) < day {
					continue // 日数超过该年该月总天数
				}

			} else {
				re := regexp.MustCompile("^([0-9]{1,2})M([0-9]{1,2})D$").FindStringSubmatch(gfK)
				if len(re) != 3 {
					continue // 索引格式不正确
				}
				// 月数
				if re[1] == "" {
					continue // 索引格式不正确,没指明月份
				}
				m = re[1]
				// 月数string转int
				month, err = strconv.Atoi(m)
				if err != nil || month < 1 || month > 12 {
					continue // 月份为空或数字不正确
				}

				// 日
				if re[2] == "" {
					continue // 索引格式不正确,没有指明日
				}
				d = re[2]
				// 日数string转int
				day, err = strconv.Atoi(d)
				if err != nil || day < 1 || day > GregorianMonthDays(gregorianYear, month) {
					continue // 日为空或数字不正确
				}

			}

			if month == 0 || day == 0 {
				continue
			}

			trueK = strconv.Itoa(month) + "M" + strconv.Itoa(day) + "D"
			// 对应值
			if _, ok := fds[trueK]; ok {
				fds[trueK] = append(fds[trueK], strings.Split(gfV, ",")...)
			} else {
				fds[trueK] = strings.Split(gfV, ",")
			}
		}

		c.tempData.gFD.setData(gregorianYear, fds)
	}

	festivalIndex := strconv.Itoa(int(gregorianMonth)) + "M" + strconv.Itoa(gregorianDay) + "D"

	var fi FestivalItem // 该日的节日
	if fv, ok := fds[festivalIndex]; ok {
		for _, v := range fv {
			v = strings.TrimSpace(v)
			if svs := strings.Split(v, "*"); len(svs) > 1 {
				fi.Show = append(fi.Show, svs[1])
			} else {
				fi.Secondary = append(fi.Secondary, v)
			}
		}
	}

	return fi
}

// (*Calendar) SolarTerms 一整年的节气
//
// 从上一年的冬至开始到下一年的小寒共26个节气对应的日期时间,
// 设置c.stMap的值 map[string]*SolarTerm是一个以"年-月-日"为索引的map,
// 返回[]*SolarTerm是一个有序切片
func (c *Calendar) SolarTerms(year int) []*SolarTermItem {
	sts := c.tempData.st.getData(year)
	if sts != nil && len(sts) == 26 {
		return sts
	}

	ji := -1

	lastYearAsts := lastYearSolarTerms(float64(year))

	for i, v := range lastYearAsts {
		if v == 0 {
			continue
		}
		if i < 18 {
			continue
		}
		if i > 23 {
			continue
		}

		ji++

		// var stItem SolarTerm
		stItem := new(SolarTermItem)

		stTime := JdToTime(v, c.loc)
		stItem.Index = (ji + 18) % 24 // 节气名称的索引
		stItem.Name = solarTermsNameArray[stItem.Index]
		stItem.Time = &stTime

		sts = append(sts, stItem)
	}

	asts := adjustedSolarTermsJd(float64(year), 0, 19)
	for i, v := range asts {

		if v == 0 {
			continue
		}

		if i > 19 {
			continue
		}

		ji++

		stItem := new(SolarTermItem)

		stTime := JdToTime(v, c.loc)
		stItem.Index = (ji + 18) % 24 // 节气名称的索引
		stItem.Name = solarTermsNameArray[stItem.Index]
		stItem.Time = &stTime

		sts = append(sts, stItem)
	}

	c.tempData.st.setData(year, sts)

	return sts
}

// StarSign 根据月和日取星座
func StarSign(month, day int) (int, string, error) {

	if month < 1 || month > 12 || day < 1 || day > 31 {
		return 0, "", errors.New("日期错误！")
	}

	// 星座的起始日期
	ZodiacDayArray := [12]int{20, 19, 21, 20, 21, 22, 23, 23, 23, 24, 22, 22}

	i := month - 1
	if day < ZodiacDayArray[i] {
		i = ((i + 12) - 1) % 12
	}

	return i, starSignsNameArray[i], nil
}

// (SolarTermItem) String 节气显示
func (sti SolarTermItem) String() string {
	return fmt.Sprintf("%s 定%s:%s", sti.Name, sti.Name, sti.Time.Format(time.RFC3339))
}

// (CalendarItem) String 日历单元显示
func (ci CalendarItem) String() string {
	dateString := ci.Time.Format("2006-01-02")
	todayStr := " "
	if ci.IsToday == 1 {
		todayStr = "*"
	}
	weekIndex := ci.Time.Weekday()
	weekString := fmt.Sprintf("%s周%s", todayStr, weekNameArray[weekIndex])

	festivalString := ""
	if ci.Festival != nil && ci.Festival.Show != nil {
		festivalString = " " + strings.Join(ci.Festival.Show, ",")
	}

	var lunarString = ""
	if ci.LunarDate != nil {
		// lunarString = " " + ci.LunarDate.String()
		// 农历节日覆盖公历节日
		if ci.LunarDate.Festival != nil && ci.LunarDate.Festival.Show != nil {
			festivalString = " " + strings.Join(ci.LunarDate.Festival.Show, ",")
		}

		lunarString = fmt.Sprintf(" %d%s%s(%s)年%s%s月%s", ci.LunarDate.Year, ci.LunarDate.YearGZ.HSN, ci.LunarDate.YearGZ.EBN, ci.LunarDate.AnimalName, ci.LunarDate.LeapStr, ci.LunarDate.MonthName, ci.LunarDate.DayName)
	}

	var gzString = ""
	if ci.GZ != nil {
		gzStr := ci.GZ.String()
		subByte := "日"
		el := strings.Index(gzStr, subByte) + len(subByte)
		gzString = " " + string([]byte(gzStr)[:el])
	}

	var solarTermString = ""
	if ci.SolarTerm != nil {
		solarTermString = " " + ci.SolarTerm.String()
	}

	return StringSplice(dateString, weekString, lunarString, gzString, solarTermString, festivalString)
}

// (*yearSolarTermTemp) getData 读节气缓存年表
func (ystt *yearSolarTermTemp) getData(k int) []*SolarTermItem {
	ystt.mu.RLock()
	defer ystt.mu.RUnlock()
	var rv []*SolarTermItem
	if ystt.data == nil {
		ystt.data = make(map[int][]*SolarTermItem)
	}
	if v, ok := ystt.data[k]; ok {
		rv = v
	}

	return rv
}

// (*yearSolarTermTemp) getData 写节气缓存年表
func (ystt *yearSolarTermTemp) setData(k int, v []*SolarTermItem) {
	ystt.mu.Lock()
	defer ystt.mu.Unlock()
	if ystt.data == nil {
		ystt.data = make(map[int][]*SolarTermItem)
	}
	ystt.data[k] = v
}

// (*yearFestivalTemp) getData 读节日缓存年表
func (yft *yearFestivalTemp) getData(k int) map[string][]string {
	yft.mu.RLock()
	defer yft.mu.RUnlock()

	rv := make(map[string][]string)
	if yft.data == nil {
		yft.data = make(map[int]map[string][]string)
	}
	if v, ok := yft.data[k]; ok {
		rv = v
	}

	return rv
}

// (*yearFestivalTemp) setData 写节日缓存年表
func (yft *yearFestivalTemp) setData(k int, v map[string][]string) {
	yft.mu.Lock()
	defer yft.mu.Unlock()
	if yft.data == nil {
		yft.data = make(map[int]map[string][]string)
	}
	yft.data[k] = v
}
