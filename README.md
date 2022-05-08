# 日历

calendar、日历、中国农历、阴历、节气、干支、生肖、星座

通过天文计算和民间推算方法，准确计算出公历-1000年至3000年的农历、干支、节气等。

> 天文计算方法参考Jean Meeus的《Astronomical Algorithms》、[NASA](https://eclipse.gsfc.nasa.gov/SEhelp/deltatpoly2004.html "NASA")网站、[天文与历法](http://www.bieyu.com/ "天文与历法")网站等相关的天文历法计算方法。

当前稳定版本(Current Stable Version)：v1.1.0

推荐版本(Recommended Version):v1.1.0

`go get github.com/liujiawm/gocalendar@v1.1.0`

- [Installation 安装](#installation-安装)
- [示例](#示例)
  - [日历表](#日历表)
  - [日历配置](#日历配置)
  - [其它接口](#其它接口)
    - [节气](#节气)
    - [农历与公历互换](#农历与公历互换)
    - [早晚子时示例说明](#早晚子时示例说明)
    - [公历转换干支](#公历转换干支)
    - [星座](#星座)
    - [儒略日(Julian Day)](#儒略日julian-day)
    - [Modified Julian Day](#modified-julian-day)
- [Documentation 更多详细说明](https://pkg.go.dev/github.com/liujiawm/gocalendar)
- [帮助](https://github.com/liujiawm/gocalendar)
- 联系
  - QQ:194088
  - Email:liujiawm@msn.com

## Installation 安装 ##

```
go get github.com/liujiawm/gocalendar
```

## 示例 ##

### 日历表 ###

#### 生成一个日历表 ####

`(*Calendar) Generate()`

``` go
// 用默认的Calendar生成日历表,当前时间是2021年2月8日
items := DefaultCalendar().Generate()
for _,item := range items {
	fmt.Println(item)
}

```

```
 公历日期(*)周        农历日期          干支           节气、节日
------------------------------------------------------------------------
2021-01-31 周日 2020庚子(鼠)年腊月十九 庚子年己丑月己卯日
2021-02-01 周一 2020庚子(鼠)年腊月二十 庚子年己丑月庚辰日
2021-02-02 周二 2020庚子(鼠)年腊月廿一 庚子年己丑月辛巳日
2021-02-03 周三 2020庚子(鼠)年腊月廿二 辛丑年庚寅月壬午日 立春 定立春:2021-02-03T22:59:23+08:00
2021-02-04 周四 2020庚子(鼠)年腊月廿三 辛丑年庚寅月癸未日
2021-02-05 周五 2020庚子(鼠)年腊月廿四 辛丑年庚寅月甲申日 小年
2021-02-06 周六 2020庚子(鼠)年腊月廿五 辛丑年庚寅月乙酉日
2021-02-07 周日 2020庚子(鼠)年腊月廿六 辛丑年庚寅月丙戌日
2021-02-08*周一 2020庚子(鼠)年腊月廿七 辛丑年庚寅月丁亥日
2021-02-09 周二 2020庚子(鼠)年腊月廿八 辛丑年庚寅月戊子日
2021-02-10 周三 2020庚子(鼠)年腊月廿九 辛丑年庚寅月己丑日
2021-02-11 周四 2020庚子(鼠)年腊月三十 辛丑年庚寅月庚寅日 除夕
2021-02-12 周五 2021辛丑(牛)年正月初一 辛丑年庚寅月辛卯日 春节
2021-02-13 周六 2021辛丑(牛)年正月初二 辛丑年庚寅月壬辰日
2021-02-14 周日 2021辛丑(牛)年正月初三 辛丑年庚寅月癸巳日 情人节
2021-02-15 周一 2021辛丑(牛)年正月初四 辛丑年庚寅月甲午日
2021-02-16 周二 2021辛丑(牛)年正月初五 辛丑年庚寅月乙未日
2021-02-17 周三 2021辛丑(牛)年正月初六 辛丑年庚寅月丙申日
2021-02-18 周四 2021辛丑(牛)年正月初七 辛丑年庚寅月丁酉日 雨水 定雨水:2021-02-18T18:44:29+08:00
2021-02-19 周五 2021辛丑(牛)年正月初八 辛丑年庚寅月戊戌日
2021-02-20 周六 2021辛丑(牛)年正月初九 辛丑年庚寅月己亥日
2021-02-21 周日 2021辛丑(牛)年正月初十 辛丑年庚寅月庚子日
2021-02-22 周一 2021辛丑(牛)年正月十一 辛丑年庚寅月辛丑日
2021-02-23 周二 2021辛丑(牛)年正月十二 辛丑年庚寅月壬寅日
2021-02-24 周三 2021辛丑(牛)年正月十三 辛丑年庚寅月癸卯日
2021-02-25 周四 2021辛丑(牛)年正月十四 辛丑年庚寅月甲辰日
2021-02-26 周五 2021辛丑(牛)年正月十五 辛丑年庚寅月乙巳日 元宵节
2021-02-27 周六 2021辛丑(牛)年正月十六 辛丑年庚寅月丙午日
2021-02-28 周日 2021辛丑(牛)年正月十七 辛丑年庚寅月丁未日
2021-03-01 周一 2021辛丑(牛)年正月十八 辛丑年庚寅月戊申日
2021-03-02 周二 2021辛丑(牛)年正月十九 辛丑年庚寅月己酉日
2021-03-03 周三 2021辛丑(牛)年正月二十 辛丑年庚寅月庚戌日
2021-03-04 周四 2021辛丑(牛)年正月廿一 辛丑年庚寅月辛亥日
2021-03-05 周五 2021辛丑(牛)年正月廿二 辛丑年辛卯月壬子日 惊蛰 定惊蛰:2021-03-05T16:53:57+08:00
2021-03-06 周六 2021辛丑(牛)年正月廿三 辛丑年辛卯月癸丑日
2021-03-07 周日 2021辛丑(牛)年正月廿四 辛丑年辛卯月甲寅日
2021-03-08 周一 2021辛丑(牛)年正月廿五 辛丑年辛卯月乙卯日 妇女节
2021-03-09 周二 2021辛丑(牛)年正月廿六 辛丑年辛卯月丙辰日
2021-03-10 周三 2021辛丑(牛)年正月廿七 辛丑年辛卯月丁巳日
2021-03-11 周四 2021辛丑(牛)年正月廿八 辛丑年辛卯月戊午日
2021-03-12 周五 2021辛丑(牛)年正月廿九 辛丑年辛卯月己未日 植树节
2021-03-13 周六 2021辛丑(牛)年二月初一 辛丑年辛卯月庚申日
```

#### 指定日期的日历表 ####

`(*Calendar) GenerateWithDate(year, month, day int, timeParts ...int)`

``` go
items := DefaultCalendar().GenerateWithDate(2021,5,1)
for _,item := range items {
	fmt.Println(item)
}
```

```
2021-04-25 周日 2021辛丑(牛)年三月十四 辛丑年壬辰月癸卯日
2021-04-26 周一 2021辛丑(牛)年三月十五 辛丑年壬辰月甲辰日
2021-04-27 周二 2021辛丑(牛)年三月十六 辛丑年壬辰月乙巳日
2021-04-28 周三 2021辛丑(牛)年三月十七 辛丑年壬辰月丙午日
2021-04-29 周四 2021辛丑(牛)年三月十八 辛丑年壬辰月丁未日
2021-04-30 周五 2021辛丑(牛)年三月十九 辛丑年壬辰月戊申日
2021-05-01 周六 2021辛丑(牛)年三月二十 辛丑年壬辰月己酉日 劳动节
2021-05-02 周日 2021辛丑(牛)年三月廿一 辛丑年壬辰月庚戌日
2021-05-03 周一 2021辛丑(牛)年三月廿二 辛丑年壬辰月辛亥日
2021-05-04 周二 2021辛丑(牛)年三月廿三 辛丑年壬辰月壬子日 青年节
2021-05-05 周三 2021辛丑(牛)年三月廿四 辛丑年癸巳月癸丑日 立夏 定立夏:2021-05-05T14:46:29+08:00
2021-05-06 周四 2021辛丑(牛)年三月廿五 辛丑年癸巳月甲寅日
2021-05-07 周五 2021辛丑(牛)年三月廿六 辛丑年癸巳月乙卯日
2021-05-08 周六 2021辛丑(牛)年三月廿七 辛丑年癸巳月丙辰日
2021-05-09 周日 2021辛丑(牛)年三月廿八 辛丑年癸巳月丁巳日 母亲节
2021-05-10 周一 2021辛丑(牛)年三月廿九 辛丑年癸巳月戊午日
2021-05-11 周二 2021辛丑(牛)年三月三十 辛丑年癸巳月己未日
2021-05-12 周三 2021辛丑(牛)年四月初一 辛丑年癸巳月庚申日 护士节
2021-05-13 周四 2021辛丑(牛)年四月初二 辛丑年癸巳月辛酉日
2021-05-14 周五 2021辛丑(牛)年四月初三 辛丑年癸巳月壬戌日
2021-05-15 周六 2021辛丑(牛)年四月初四 辛丑年癸巳月癸亥日
2021-05-16 周日 2021辛丑(牛)年四月初五 辛丑年癸巳月甲子日
2021-05-17 周一 2021辛丑(牛)年四月初六 辛丑年癸巳月乙丑日
2021-05-18 周二 2021辛丑(牛)年四月初七 辛丑年癸巳月丙寅日
2021-05-19 周三 2021辛丑(牛)年四月初八 辛丑年癸巳月丁卯日
2021-05-20 周四 2021辛丑(牛)年四月初九 辛丑年癸巳月戊辰日
2021-05-21 周五 2021辛丑(牛)年四月初十 辛丑年癸巳月己巳日 小满 定小满:2021-05-21T03:36:22+08:00
2021-05-22 周六 2021辛丑(牛)年四月十一 辛丑年癸巳月庚午日
2021-05-23 周日 2021辛丑(牛)年四月十二 辛丑年癸巳月辛未日
2021-05-24 周一 2021辛丑(牛)年四月十三 辛丑年癸巳月壬申日
2021-05-25 周二 2021辛丑(牛)年四月十四 辛丑年癸巳月癸酉日
2021-05-26 周三 2021辛丑(牛)年四月十五 辛丑年癸巳月甲戌日
2021-05-27 周四 2021辛丑(牛)年四月十六 辛丑年癸巳月乙亥日
2021-05-28 周五 2021辛丑(牛)年四月十七 辛丑年癸巳月丙子日
2021-05-29 周六 2021辛丑(牛)年四月十八 辛丑年癸巳月丁丑日
2021-05-30 周日 2021辛丑(牛)年四月十九 辛丑年癸巳月戊寅日
2021-05-31 周一 2021辛丑(牛)年四月二十 辛丑年癸巳月己卯日
2021-06-01 周二 2021辛丑(牛)年四月廿一 辛丑年癸巳月庚辰日 儿童节
2021-06-02 周三 2021辛丑(牛)年四月廿二 辛丑年癸巳月辛巳日
2021-06-03 周四 2021辛丑(牛)年四月廿三 辛丑年癸巳月壬午日
2021-06-04 周五 2021辛丑(牛)年四月廿四 辛丑年癸巳月癸未日
2021-06-05 周六 2021辛丑(牛)年四月廿五 辛丑年甲午月甲申日 芒种 定芒种:2021-06-05T18:51:32+08:00
```

#### 日历变换 ####

下一月

`(*Calendar) NextMonth`

``` go
c := DefaultCalendar()

// 下一月的日历表
items := c.NextMonth()

```

上一月

`(*Calendar) PreviousMonth`

``` go
c := DefaultCalendar()

// 上一月的日历表
items := c.PreviousMonth()

```

下一年当月

`(*Calendar) NextYear`

``` go
c := DefaultCalendar()

// 下一年当月的日历表
items := c.NextYear()

```

上一年当月

`(*Calendar) PreviousYear`

``` go
c := DefaultCalendar()

// 上一年的当月日历表
items := c.PreviousYear()

```

### 日历配置 ###

``` go
type CalendarConfig struct {
	Grid            int    // 取日历方式,GridDay按天取日历,GridWeek按周取日历,GridMonth按月取日历
	FirstWeek       int    // 日历显示时第一列显示周几，(日历表第一列是周几,0周日,依次最大值6)
	TimeZoneName    string // 时区名称,需zoneinfo支持的时区名称
	SolarTerms      bool   // 读取节气 bool
	Lunar           bool   // 读取农历 bool
	HeavenlyEarthly bool   // 读取干支 bool
	NightZiHour     bool   // 区分早晚子时，true则 23:00-24:00 00:00-01:00为子时，否则00:00-02:00为子时
	StarSign        bool   // 读取星座
}

```

#### 自定义日历 ####

`NewCalendar(CalendarConfig)`

``` go
c := NewCalendar(CalendarConfig{
		Grid:GridWeek,
		FirstWeek:0,
		SolarTerms:true,
		Lunar:true,
		HeavenlyEarthly:true,
		NightZiHour:true,
		StarSign:true,
	})

items :=c.GenerateWithDate(2021,12,22)
for _,item := range result {
	fmt.Println(item)
}

```

```
2021-12-19 周日 2021辛丑(牛)年十一月十六 辛丑年庚子月辛丑日
2021-12-20 周一 2021辛丑(牛)年十一月十七 辛丑年庚子月壬寅日
2021-12-21 周二 2021辛丑(牛)年十一月十八 辛丑年庚子月癸卯日 冬至 定冬至:2021-12-21T23:59:05+08:00
2021-12-22 周三 2021辛丑(牛)年十一月十九 辛丑年庚子月甲辰日
2021-12-23 周四 2021辛丑(牛)年十一月二十 辛丑年庚子月乙巳日
2021-12-24 周五 2021辛丑(牛)年十一月廿一 辛丑年庚子月丙午日
2021-12-25 周六 2021辛丑(牛)年十一月廿二 辛丑年庚子月丁未日 圣诞节
```

### 其它接口 ###

#### 节气 ####

> 注意：节气依春分点计算，不同时区因时差不同，定节气时间也会不同

全年节气

`(*Calendar) SolarTerms(year int) []*SolarTerm`

``` go
// c := DefaultCalendar()

// 该package默认为本地时区，如自定议时区 Asia/Shanghai ,这将按中国时区计算节气
c := NewCalendar(CalendarConfig{TimeZoneName:"Asia/Shanghai"})
sts:= c.SolarTerms(2021)
fmt.Println("2021年节气:")
for _,v := range sts{
	fmt.Printf(" %s 定%s:%s \n", v.Name, v.Name, v.Time.Format(time.RFC3339))
}

```

```
2021年节气:
 冬至 定冬至:2020-12-21T18:02:36+08:00 
 小寒 定小寒:2021-01-05T11:23:50+08:00 
 大寒 定大寒:2021-01-20T04:40:31+08:00 
 立春 定立春:2021-02-03T22:59:23+08:00 
 雨水 定雨水:2021-02-18T18:44:29+08:00 
 惊蛰 定惊蛰:2021-03-05T16:53:57+08:00 
 春分 定春分:2021-03-20T17:37:28+08:00 
 清明 定清明:2021-04-04T21:34:48+08:00 
 谷雨 定谷雨:2021-04-20T04:32:43+08:00 
 立夏 定立夏:2021-05-05T14:46:29+08:00 
 小满 定小满:2021-05-21T03:36:22+08:00 
 芒种 定芒种:2021-06-05T18:51:32+08:00 
 夏至 定夏至:2021-06-21T11:31:47+08:00 
 小暑 定小暑:2021-07-07T05:05:28+08:00 
 大暑 定大暑:2021-07-22T22:26:42+08:00 
 立秋 定立秋:2021-08-07T14:54:28+08:00 
 处暑 定处暑:2021-08-23T05:35:23+08:00 
 白露 定白露:2021-09-07T17:53:16+08:00 
 秋分 定秋分:2021-09-23T03:20:56+08:00 
 寒露 定寒露:2021-10-08T09:38:45+08:00 
 霜降 定霜降:2021-10-23T12:50:30+08:00 
 立冬 定立冬:2021-11-07T12:58:14+08:00 
 小雪 定小雪:2021-11-22T10:33:05+08:00 
 大雪 定大雪:2021-12-07T05:56:49+08:00 
 冬至 定冬至:2021-12-21T23:59:05+08:00 
 小寒 定小寒:2022-01-05T17:14:07+08:00 
```

#### 农历与公历互换 ####

> 农历以中国(东八区)时区对应公历

公历转农历

`(*Calendar) GregorianToLunar(year, month, day int) LunarDate`

``` go
ld := DefaultCalendar().GregorianToLunar(2020,6,5)
fmt.Printf("%d%s%s(%s)年%s%s月%s\n",ld.Year,ld.YearGZ.HSN,ld.YearGZ.EBN,ld.AnimalName,ld.LeapStr,ld.MonthName,ld.DayName)

```

```
2020庚子(鼠)年闰四月十四
```

农历转公历

`(*Calendar) LunarToGregorian(lunarYear,lunarMonth,lunarDay int, isLeap bool) (time.Time, error)`

``` go
c := DefaultCalendar()
gd,_ := c.LunarToGregorian(2020,4,14,false)
fmt.Println("农历2020年四月十四转换成公历是:", gd.Format("2006-01-02"))

gd,_ = c.LunarToGregorian(2020,4,14,true)
fmt.Println("农历2020年闰四月十四转换成公历是:", gd.Format("2006-01-02"))

```

```
农历2020年四月十四转换成公历是: 2020-05-06
农历2020年闰四月十四转换成公历是: 2020-06-05
```

#### 早晚子时示例说明 ####

干支四柱的子时是23:00-00:00 00:00-01:00 说明每日开始是从上一日的23点开始

我们在该日历中引进早晚子时，允许把子时分成晚子时和早子时

如果不区分早晚子时，每日将依23点划分，如果区分早晚子时，则依0点划分

看早晚子时区分的结果

``` go
// NightZiHour:false 不区分早晚子时
c := NewCalendar(CalendarConfig{NightZiHour:false})
rt := time.Date(2021,5,6,23,50,0,0,time.Local)
gz := c.ChineseSexagenaryCycle(rt)
fmt.Printf("%s%s年%s%s月(%s%s日)%s%s时\n", gz.Year.HSN, gz.Year.EBN, gz.Month.HSN, gz.Month.EBN, gz.Day.HSN, gz.Day.EBN, gz.Hour.HSN, gz.Hour.EBN)

// NightZiHour:true 区分早晚子时
c = NewCalendar(CalendarConfig{NightZiHour:true})
rt = time.Date(2021,5,6,23,50,0,0,time.Local)
gz = c.ChineseSexagenaryCycle(rt)
fmt.Printf("%s%s年%s%s月(%s%s日)%s%s时\n", gz.Year.HSN, gz.Year.EBN, gz.Month.HSN, gz.Month.EBN, gz.Day.HSN, gz.Day.EBN, gz.Hour.HSN, gz.Hour.EBN)

```

> 结果中括号内有不同，该结果只在23点至00点会有所表现

```
辛丑年癸巳月(乙卯日)丙子时
辛丑年癸巳月(甲寅日)丙子时
```

#### 公历转换干支 ####

> 四柱干支以中国(东八区)时区对应公历

`(*Calendar) ChineseSexagenaryCycle(time.Time)GZ{}` 日期时间对应的干支

``` go
rt := time.Date(2021,5,6,23,50,0,0,time.Local)
gz := DefaultCalendar().ChineseSexagenaryCycle(rt)
fmt.Printf("%s%s年%s%s月%s%s日%s%s时\n", gz.Year.HSN, gz.Year.EBN, gz.Month.HSN, gz.Month.EBN, gz.Day.HSN, gz.Day.EBN, gz.Hour.HSN, gz.Hour.EBN)
```

```
辛丑年癸巳月甲寅日丙子时
```

#### 星座 ####

`StarSign(month,day int)(int, string, error)`

``` go
// i,ss,err := StarSign(5,6)

i,ss,_ := StarSign(5,6)
fmt.Println(ss)
// 金牛
```

#### 儒略日(Julian Day) ####

日期时间转儒略日

`JulianDay(year, month, day float64, timeParts ...float64) float64`

``` go
jd := JulianDay(2021,12,6)
// 2459554.5

jd = JulianDay(2021,12,6,12)
// 2459555

jd = JulianDay(2021,12,6,12,10,10)
// 2459555.007060185
```

儒略日转日期时间

方法一:`JdToTimeMap(jd float64) map[string]int` 返回一个日期时间map

方法二:`JdToTime(jd float64, loc *time.Location) time.Time` 返回time.Time

``` go
datetime := JdToTimeMap(2459554.5)
fmt.Printf("%d年%d月%d日\n",datetime["year"], datetime["month"],datetime["day"])
// 2021年12月6日

datetime = JdToTimeMap(2459555)
fmt.Printf("%d年%d月%d日%d时\n",datetime["year"], datetime["month"],datetime["day"],datetime["hour"])
// 2021年12月6日12时

datetime = JdToTimeMap(2459555.007060185)
fmt.Printf("%d年%d月%d日%d时%d分%d秒\n",datetime["year"], datetime["month"],datetime["day"],datetime["hour"],datetime["minute"],datetime["second"])
// 2021年12月6日12时10分10秒

// 方法二
// datetime := JdToTime(2459555.007060185,nil)
datetime := JdToTime(2459555.007060185,time.Local)
fmt.Println(datetime.Format(time.RFC3339))
// 2021-12-06T20:10:10+08:00
```

#### Modified Julian Day ####

`Mjd(year, month, day float64, timeParts ...float64) float64`

``` go
mjd := Mjd(2021,12,6)
// 59554

mjd = Mjd(2021,12,6,12)
// 59554.5

mjd = Mjd(2021,12,6,12,10,10)
// 59554.5070601851

```

`MjdToTimeMap(mjd float64) map[string]int`

``` go
datetime := MjdToTimeMap(59554)
fmt.Printf("%d年%d月%d日\n",datetime["year"], datetime["month"],datetime["day"])
// 2021年12月6日

datetime = MjdToTimeMap(59554.5)
fmt.Printf("%d年%d月%d日%d时\n",datetime["year"], datetime["month"],datetime["day"],datetime["hour"])
// 2021年12月6日12时

datetime = MjdToTimeMap(59554.507060185075)
fmt.Printf("%d年%d月%d日%d时%d分%d秒\n",datetime["year"], datetime["month"],datetime["day"],datetime["hour"],datetime["minute"],datetime["second"])
// 2021年12月6日12时10分10秒
```
