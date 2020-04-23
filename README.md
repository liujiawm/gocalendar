# gocalendar

一个用golang写的日历，有公历转农历，农历转公历，节气，干支，星座，生肖等功能

中国的农历历法综合了太阳历和月亮历,为中国的生活生产提供了重要的帮助,是中国古人智慧与中国传统文化的一个重要体现

程序比较准确的计算出农历与二十四节气(精确到分),时间限制在1000-3000年间,在实际生产使用中注意限制年份

Author: liujiawm@gmail.com

Version: "1.0.2"

日历按月取时固定为42个日期，按周取时固定为7个日期

# 测试代码

## 日历:公历和农历
```
func BenchmarkCalendar_Calendars(b *testing.B) {
	sc := Calendar{
		Loc:       time.Local, // 时区,默认time.Local
		FirstWeek: 1,          // 日历显示时第一列显示周几，(日历表第一列是周几,0周日,依次最大值6)
		Grid:      GridMonth,  // 取日历方式,默认取一个月, gocalendar.GridDay取一天,gocalendar.GridWeek 取一周,gocalendar.GridMonth取一个月
		Zwz:       false,      // 是否区分早晚子时(子时从23-01时),true则23:00-24:00算成上一天
		Getjq:     true,       // 是否取节气
	}

	cds := sc.Calendars(2020,6,5)
	for _,cd := range cds {
		var jqStr string
		if cd.SD.Jq != nil {
			jqStr = fmt.Sprintf(" 节气:%s (定%s：%d:%d:%d)",cd.SD.Jq.Name,cd.SD.Jq.Name,cd.SD.Jq.Date.Hour,cd.SD.Jq.Date.Min,cd.SD.Jq.Date.Sec)
		}

		b.Logf("公历:%d年%d月%d日 周%s 农历:%d(%s%s)[%s]年%s%s月%s (%s%s年%s%s月%s%s日%s%s时) %s",
			cd.SD.Date.Year, cd.SD.Date.Month, cd.SD.Date.Day, WeekChinese(cd.SD.Date.Week),
			cd.LD.Date.Year,cd.LD.YearGanZi.Gan,cd.LD.YearGanZi.Zhi,cd.LD.YearGanZi.Animal, cd.LD.LeapStr, cd.LD.MonthStr, cd.LD.DayStr,
			cd.SD.GanZhi.YtgStr, cd.SD.GanZhi.YdzStr, cd.SD.GanZhi.MtgStr, cd.SD.GanZhi.MdzStr, cd.SD.GanZhi.DtgStr, cd.SD.GanZhi.DdzStr, cd.SD.GanZhi.HtgStr, cd.SD.GanZhi.HdzStr,
			jqStr)
	}
}

```

公历和农历测试结果显示如下：

```
goos: windows
goarch: amd64
pkg: github.com/liujiawm/gocalendar
BenchmarkCalendar_Calendars-4   	1000000000	         0.00252 ns/op
--- BENCH: BenchmarkCalendar_Calendars-4
    calendar_test.go:44: 公历:2020年6月1日 周一 农历:2020(庚子)[鼠]年闰四月初十 (庚子年辛巳月乙亥日丁亥时) 
    calendar_test.go:44: 公历:2020年6月2日 周二 农历:2020(庚子)[鼠]年闰四月十一 (庚子年辛巳月丙子日己亥时) 
    calendar_test.go:44: 公历:2020年6月3日 周三 农历:2020(庚子)[鼠]年闰四月十二 (庚子年辛巳月丁丑日辛亥时) 
    calendar_test.go:44: 公历:2020年6月4日 周四 农历:2020(庚子)[鼠]年闰四月十三 (庚子年辛巳月戊寅日癸亥时) 
    calendar_test.go:44: 公历:2020年6月5日 周五 农历:2020(庚子)[鼠]年闰四月十四 (庚子年壬午月己卯日乙亥时)  节气:芒种 (定芒种：12:57:52)
    calendar_test.go:44: 公历:2020年6月6日 周六 农历:2020(庚子)[鼠]年闰四月十五 (庚子年壬午月庚辰日丁亥时) 
    calendar_test.go:44: 公历:2020年6月7日 周日 农历:2020(庚子)[鼠]年闰四月十六 (庚子年壬午月辛巳日己亥时) 
    calendar_test.go:44: 公历:2020年6月8日 周一 农历:2020(庚子)[鼠]年闰四月十七 (庚子年壬午月壬午日辛亥时) 
    calendar_test.go:44: 公历:2020年6月9日 周二 农历:2020(庚子)[鼠]年闰四月十八 (庚子年壬午月癸未日癸亥时) 
    calendar_test.go:44: 公历:2020年6月10日 周三 农历:2020(庚子)[鼠]年闰四月十九 (庚子年壬午月甲申日乙亥时) 
	... [output truncated]
PASS

```

## 日历:公历
```
func TestCalendar_SolarCalendar(t *testing.T) {
	// 默认设置 DefaultCalendar()
	datas2 := DefaultCalendar().SolarCalendar(2020,12,31)
	if datas2 != nil {
		for _,d := range datas2 {
			var jqstr string
			if d.Jq != nil {
				jqstr = fmt.Sprintf(" 节气:%s (定%s：%d:%d:%d) \n",d.Jq.Name,d.Jq.Name,d.Jq.Date.Hour,d.Jq.Date.Min,d.Jq.Date.Sec)
			}
			t.Logf("%d年%d月%d日 %d时%d分%d秒 周%d (%s%s年%s%s月%s%s日%s%s时) %s",d.Year,d.Month,d.Day,d.Hour,d.Min,d.Sec,d.Week,
				d.GanZhi.YtgStr, d.GanZhi.YdzStr, d.GanZhi.MtgStr, d.GanZhi.MdzStr, d.GanZhi.DtgStr, d.GanZhi.DdzStr, d.GanZhi.HtgStr, d.GanZhi.HdzStr,
				jqstr)
		}
	}
}

```

公历测试结果显示如下:

```
=== RUN   TestCalendar_SolarCalendar
--- PASS: TestCalendar_SolarCalendar (0.02s)
    calendar_test.go:61: 2020年11月29日 21时14分39秒 周0 (庚子年丁亥月丙子日己亥时) 
    calendar_test.go:61: 2020年11月30日 21时14分39秒 周1 (庚子年丁亥月丁丑日辛亥时) 
    calendar_test.go:61: 2020年12月1日 21时14分39秒 周2 (庚子年丁亥月戊寅日癸亥时) 
    calendar_test.go:61: 2020年12月2日 21时14分39秒 周3 (庚子年丁亥月己卯日乙亥时) 
    calendar_test.go:61: 2020年12月3日 21时14分39秒 周4 (庚子年丁亥月庚辰日丁亥时) 
    calendar_test.go:61: 2020年12月4日 21时14分39秒 周5 (庚子年丁亥月辛巳日己亥时) 
    calendar_test.go:61: 2020年12月5日 21时14分39秒 周6 (庚子年丁亥月壬午日辛亥时) 
    calendar_test.go:61: 2020年12月6日 21时14分39秒 周0 (庚子年丁亥月癸未日癸亥时) 
    calendar_test.go:61: 2020年12月7日 21时14分39秒 周1 (庚子年戊子月甲申日乙亥时)  节气:大雪 (定大雪：0:9:23) 
    calendar_test.go:61: 2020年12月8日 21时14分39秒 周2 (庚子年戊子月乙酉日丁亥时) 
    calendar_test.go:61: 2020年12月9日 21时14分39秒 周3 (庚子年戊子月丙戌日己亥时) 
    calendar_test.go:61: 2020年12月10日 21时14分39秒 周4 (庚子年戊子月丁亥日辛亥时) 
    calendar_test.go:61: 2020年12月11日 21时14分39秒 周5 (庚子年戊子月戊子日癸亥时) 
    calendar_test.go:61: 2020年12月12日 21时14分39秒 周6 (庚子年戊子月己丑日乙亥时) 
    calendar_test.go:61: 2020年12月13日 21时14分39秒 周0 (庚子年戊子月庚寅日丁亥时) 
    calendar_test.go:61: 2020年12月14日 21时14分39秒 周1 (庚子年戊子月辛卯日己亥时) 
    calendar_test.go:61: 2020年12月15日 21时14分39秒 周2 (庚子年戊子月壬辰日辛亥时) 
    calendar_test.go:61: 2020年12月16日 21时14分39秒 周3 (庚子年戊子月癸巳日癸亥时) 
    calendar_test.go:61: 2020年12月17日 21时14分39秒 周4 (庚子年戊子月甲午日乙亥时) 
    calendar_test.go:61: 2020年12月18日 21时14分39秒 周5 (庚子年戊子月乙未日丁亥时) 
    calendar_test.go:61: 2020年12月19日 21时14分39秒 周6 (庚子年戊子月丙申日己亥时) 
    calendar_test.go:61: 2020年12月20日 21时14分39秒 周0 (庚子年戊子月丁酉日辛亥时) 
    calendar_test.go:61: 2020年12月21日 21时14分39秒 周1 (庚子年戊子月戊戌日癸亥时)  节气:冬至 (定冬至：18:2:36) 
    calendar_test.go:61: 2020年12月22日 21时14分39秒 周2 (庚子年戊子月己亥日乙亥时) 
    calendar_test.go:61: 2020年12月23日 21时14分39秒 周3 (庚子年戊子月庚子日丁亥时) 
    calendar_test.go:61: 2020年12月24日 21时14分39秒 周4 (庚子年戊子月辛丑日己亥时) 
    calendar_test.go:61: 2020年12月25日 21时14分39秒 周5 (庚子年戊子月壬寅日辛亥时) 
    calendar_test.go:61: 2020年12月26日 21时14分39秒 周6 (庚子年戊子月癸卯日癸亥时) 
    calendar_test.go:61: 2020年12月27日 21时14分39秒 周0 (庚子年戊子月甲辰日乙亥时) 
    calendar_test.go:61: 2020年12月28日 21时14分39秒 周1 (庚子年戊子月乙巳日丁亥时) 
    calendar_test.go:61: 2020年12月29日 21时14分39秒 周2 (庚子年戊子月丙午日己亥时) 
    calendar_test.go:61: 2020年12月30日 21时14分39秒 周3 (庚子年戊子月丁未日辛亥时) 
    calendar_test.go:61: 2020年12月31日 21时14分39秒 周4 (庚子年戊子月戊申日癸亥时) 
    calendar_test.go:61: 2021年1月1日 21时14分39秒 周5 (庚子年戊子月己酉日乙亥时) 
    calendar_test.go:61: 2021年1月2日 21时14分39秒 周6 (庚子年戊子月庚戌日丁亥时) 
    calendar_test.go:61: 2021年1月3日 21时14分39秒 周0 (庚子年戊子月辛亥日己亥时) 
    calendar_test.go:61: 2021年1月4日 21时14分39秒 周1 (庚子年戊子月壬子日辛亥时) 
    calendar_test.go:61: 2021年1月5日 21时14分39秒 周2 (庚子年己丑月癸丑日癸亥时)  节气:小寒 (定小寒：11:23:50) 
    calendar_test.go:61: 2021年1月6日 21时14分39秒 周3 (庚子年己丑月甲寅日乙亥时) 
    calendar_test.go:61: 2021年1月7日 21时14分39秒 周4 (庚子年己丑月乙卯日丁亥时) 
    calendar_test.go:61: 2021年1月8日 21时14分39秒 周5 (庚子年己丑月丙辰日己亥时) 
    calendar_test.go:61: 2021年1月9日 21时14分39秒 周6 (庚子年己丑月丁巳日辛亥时) 
PASS

```

# JSON

```
	sc := Calendar{
		Loc:       time.Local, // 时区,默认time.Local
		FirstWeek: 0,          // 从周几天始,默认0周日开始, (日历表第列是周几,0周日,依次最大值6)
		Grid:      GridWeek,  // 取日历方式,默认取一个月, gocalendar.GridDay取一天,gocalendar.GridWeek 取一周,gocalendar.GridMonth取一个月
		Zwz:       false,      // 是否区分早晚子时(子时从23-01时),true则23:00-24:00算成上一天
		Getjq:     true,       // 是否取节气
	}
	cds := sc.Calendars(2020,4,22)
	cdjson,err := json.Marshal(cds)
	if err != nil {
		t.Error("JSON ERR:", err.Error())
	}
	t.Logf(string(cdjson))

```

结果:

```

[{"solar":{"date":{"year":2020,"month":4,"day":19,"hour":22,"min":1,"sec":22,"week":0},"jq":{"name":"谷雨","date":{"year":2020,"month":4,"day":19,"hour":22,"min":45,"sec":11,"week":0}},"gan_zhi":{"ytg":6,"ytg_str":"庚","ydz":0,"ydz_str":"子","mtg":6,"mtg_str":"庚","mdz":4,"mdz_str":"辰","dtg":8,"dtg_str":"壬","ddz":4,"ddz_str":"辰","htg":7,"htg_str":"辛","hdz":11,"hdz_str":"亥"}},"lunar":{"date":{"year":2020,"month":3,"day":27,"hour":22,"min":1,"sec":22,"week":0},"month_str":"三","day_str":"廿七","leap_str":"","leap_year":4,"leap_month":0,"gan_zi":{"gan":"庚","zhi":"子","animal":"鼠"}}},{"solar":{"date":{"year":2020,"month":4,"day":20,"hour":22,"min":1,"sec":22,"week":1},"jq":null,"gan_zhi":{"ytg":6,"ytg_str":"庚","ydz":0,"ydz_str":"子","mtg":6,"mtg_str":"庚","mdz":4,"mdz_str":"辰","dtg":9,"dtg_str":"癸","ddz":5,"ddz_str":"巳","htg":9,"htg_str":"癸","hdz":11,"hdz_str":"亥"}},"lunar":{"date":{"year":2020,"month":3,"day":28,"hour":22,"min":1,"sec":22,"week":1},"month_str":"三","day_str":"廿八","leap_str":"","leap_year":4,"leap_month":0,"gan_zi":{"gan":"庚","zhi":"子","animal":"鼠"}}},{"solar":{"date":{"year":2020,"month":4,"day":21,"hour":22,"min":1,"sec":22,"week":2},"jq":null,"gan_zhi":{"ytg":6,"ytg_str":"庚","ydz":0,"ydz_str":"子","mtg":6,"mtg_str":"庚","mdz":4,"mdz_str":"辰","dtg":0,"dtg_str":"甲","ddz":6,"ddz_str":"午","htg":1,"htg_str":"乙","hdz":11,"hdz_str":"亥"}},"lunar":{"date":{"year":2020,"month":3,"day":29,"hour":22,"min":1,"sec":22,"week":2},"month_str":"三","day_str":"廿九","leap_str":"","leap_year":4,"leap_month":0,"gan_zi":{"gan":"庚","zhi":"子","animal":"鼠"}}},{"solar":{"date":{"year":2020,"month":4,"day":22,"hour":22,"min":1,"sec":22,"week":3},"jq":null,"gan_zhi":{"ytg":6,"ytg_str":"庚","ydz":0,"ydz_str":"子","mtg":6,"mtg_str":"庚","mdz":4,"mdz_str":"辰","dtg":1,"dtg_str":"乙","ddz":7,"ddz_str":"未","htg":3,"htg_str":"丁","hdz":11,"hdz_str":"亥"}},"lunar":{"date":{"year":2020,"month":3,"day":30,"hour":22,"min":1,"sec":22,"week":3},"month_str":"三","day_str":"卅十","leap_str":"","leap_year":4,"leap_month":0,"gan_zi":{"gan":"庚","zhi":"子","animal":"鼠"}}},{"solar":{"date":{"year":2020,"month":4,"day":23,"hour":22,"min":1,"sec":22,"week":4},"jq":null,"gan_zhi":{"ytg":6,"ytg_str":"庚","ydz":0,"ydz_str":"子","mtg":6,"mtg_str":"庚","mdz":4,"mdz_str":"辰","dtg":2,"dtg_str":"丙","ddz":8,"ddz_str":"申","htg":5,"htg_str":"己","hdz":11,"hdz_str":"亥"}},"lunar":{"date":{"year":2020,"month":4,"day":1,"hour":22,"min":1,"sec":22,"week":4},"month_str":"四","day_str":"初一","leap_str":"","leap_year":4,"leap_month":0,"gan_zi":{"gan":"庚","zhi":"子","animal":"鼠"}}},{"solar":{"date":{"year":2020,"month":4,"day":24,"hour":22,"min":1,"sec":22,"week":5},"jq":null,"gan_zhi":{"ytg":6,"ytg_str":"庚","ydz":0,"ydz_str":"子","mtg":6,"mtg_str":"庚","mdz":4,"mdz_str":"辰","dtg":3,"dtg_str":"丁","ddz":9,"ddz_str":"酉","htg":7,"htg_str":"辛","hdz":11,"hdz_str":"亥"}},"lunar":{"date":{"year":2020,"month":4,"day":2,"hour":22,"min":1,"sec":22,"week":5},"month_str":"四","day_str":"初二","leap_str":"","leap_year":4,"leap_month":0,"gan_zi":{"gan":"庚","zhi":"子","animal":"鼠"}}},{"solar":{"date":{"year":2020,"month":4,"day":25,"hour":22,"min":1,"sec":22,"week":6},"jq":null,"gan_zhi":{"ytg":6,"ytg_str":"庚","ydz":0,"ydz_str":"子","mtg":6,"mtg_str":"庚","mdz":4,"mdz_str":"辰","dtg":4,"dtg_str":"戊","ddz":10,"ddz_str":"戌","htg":9,"htg_str":"癸","hdz":11,"hdz_str":"亥"}},"lunar":{"date":{"year":2020,"month":4,"day":3,"hour":22,"min":1,"sec":22,"week":6},"month_str":"四","day_str":"初三","leap_str":"","leap_year":4,"leap_month":0,"gan_zi":{"gan":"庚","zhi":"子","animal":"鼠"}}}]

```


# 取节气

```
func TestCalendar_Jieqi(t *testing.T) {
	jq,_,err := DefaultCalendar().Jieqi(2020)
	if err != nil {
		t.Error(err.Error())
	}
	for _,d := range jq {
		t.Logf("%s - %d年%d月%d日 %d时%d分%d秒 周%d",d.Name,d.Date.Year,d.Date.Month,d.Date.Day,d.Date.Hour,d.Date.Min,d.Date.Sec,d.Date.Week)
	}
}

```

节气结果如下:

```

=== RUN   TestCalendar_Jieqi
--- PASS: TestCalendar_Jieqi (0.01s)
    calendar_test.go:86: 冬至 - 2019年12月22日 12时19分6秒 周0
    calendar_test.go:86: 小寒 - 2020年1月6日 5时30分6秒 周1
    calendar_test.go:86: 大寒 - 2020年1月20日 22时54分52秒 周1
    calendar_test.go:86: 立春 - 2020年2月4日 17时3分44秒 周2
    calendar_test.go:86: 雨水 - 2020年2月19日 12时57分27秒 周3
    calendar_test.go:86: 惊蛰 - 2020年3月5日 10时57分22秒 周4
    calendar_test.go:86: 春分 - 2020年3月20日 11时49分57秒 周5
    calendar_test.go:86: 清明 - 2020年4月4日 15时38分19秒 周6
    calendar_test.go:86: 谷雨 - 2020年4月19日 22时45分11秒 周0
    calendar_test.go:86: 立夏 - 2020年5月5日 8时50分58秒 周2
    calendar_test.go:86: 小满 - 2020年5月20日 21时48分35秒 周3
    calendar_test.go:86: 芒种 - 2020年6月5日 12时57分52秒 周5
    calendar_test.go:86: 夏至 - 2020年6月21日 5时43分16秒 周0
    calendar_test.go:86: 小暑 - 2020年7月6日 23时14分34秒 周1
    calendar_test.go:86: 大暑 - 2020年7月22日 16时37分21秒 周3
    calendar_test.go:86: 立秋 - 2020年8月7日 9时6分52秒 周5
    calendar_test.go:86: 处暑 - 2020年8月22日 23时45分38秒 周6
    calendar_test.go:86: 白露 - 2020年9月7日 12时8分26秒 周1
    calendar_test.go:86: 秋分 - 2020年9月22日 21时30分52秒 周2
    calendar_test.go:86: 寒露 - 2020年10月8日 3时55分2秒 周4
    calendar_test.go:86: 霜降 - 2020年10月23日 6时59分18秒 周5
    calendar_test.go:86: 立冬 - 2020年11月7日 7时13分27秒 周6
    calendar_test.go:86: 小雪 - 2020年11月22日 4时39分32秒 周0
    calendar_test.go:86: 大雪 - 2020年12月7日 0时9分23秒 周1
    calendar_test.go:86: 冬至 - 2020年12月21日 18时2分36秒 周1
    calendar_test.go:86: 小寒 - 2021年1月5日 11时23分50秒 周2
    calendar_test.go:86: 大寒 - 2021年1月20日 4时40分31秒 周3
PASS

```

# 公历农历互换

```

func TestCalendar_Solar2Lunar(t *testing.T) {
	sc := Calendar{
		Getjq:     false,       // 是否取节气
	}

	d := YmdNewDate(2020,5,28,time.Local)
	ld,err := sc.Solar2Lunar(sc.DateToSolarDate(d))
	if err != nil {
		t.Error(err)
	}

	t.Logf("公历%d年%d月%d日转成农历是: %d(%s%s)[%s]年%s%s月%s",d.Year,d.Month,d.Day,ld.Year,ld.YearGanZi.Gan,ld.YearGanZi.Zhi,ld.YearGanZi.Animal, ld.LeapStr, ld.MonthStr, ld.DayStr)

	ssd,err := sc.Lunar2Solar(ld.Year,ld.Month,ld.Day,ld.LeapMonth)
	if err != nil {
		t.Error(err)
	}
	t.Logf("农历%d年%d月%d日转成公历是: %d年%d月%d日",ld.Year,ld.Month,ld.Day,ssd.Year,ssd.Month,ssd.Day)

}


```

结果:

```
=== RUN   TestCalendar_Solar2Lunar
--- PASS: TestCalendar_Solar2Lunar (0.02s)
    calendar_test.go:93: 公历2020年5月28日转成农历是: 2020(庚子)[鼠]年闰四月初六
    calendar_test.go:99: 农历2020年4月6日转成公历是: 2020年5月28日
PASS

```

# 星座

```
func TestZodiac(t *testing.T) {
	// 用日期的月和日取出对应的星座索引和名称
	i,zo := Zodiac(4,22)
	t.Log(i)
	t.Log(zo)
}

```

结果

```
--- PASS: TestZodiac (0.00s)
    calendar_test.go:107: 3
    calendar_test.go:108: 金牛
PASS

```
