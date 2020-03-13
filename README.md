# gocalendar
golang写的一个日历，有公历转农历，农历转公历，节气，干支，星座，生肖等功能

日历显示查看calendar_test.go的TestGetMonthCalendar

显示42个日期，同时显示农历和节气

测试代码
```
c := GetMonthCalendar(2020,1)

	t.Log(c.Year,"年",c.Month,"月")
	t.Log("今天：",c.ToDay)
	for _,v := range c.Items {
		leapMonthStr := ""
		if v.LunarDate.LeapMonth == 1 {
			leapMonthStr = "(闰)"
		}
		t.Logf("公历年:%d年%d月%d日周%d 农历年:%d(%s%s)[%s]年%s%s月%s日 节气:%s",v.Date.Year,v.Date.Month,v.Date.Day,v.Date.Week,
			v.LunarDate.Year, v.LunarDate.YearGanZi.Gan, v.LunarDate.YearGanZi.Zhi, v.LunarDate.YearGanZi.Animals,
			leapMonthStr,MonthChinese(v.LunarDate.Month),DayChinese(v.LunarDate.Day),v.Date.JQ.Name)
	}
```

测试结果显示如下：（实际使用可转json）


```
=== RUN   TestGetMonthCalendar
--- PASS: TestGetMonthCalendar (0.01s)
    calendar_test.go:10: 2020 年 1 月
    calendar_test.go:11: 今天： &{2020 3 14 4 0 6 0 6 Local 0xc0000b9b80 0xc0000b73a0}
    calendar_test.go:17: 公历年:2019年12月29日周0 农历年:2019(己亥)[猪]年腊月初四日 节气:
    calendar_test.go:17: 公历年:2019年12月30日周1 农历年:2019(己亥)[猪]年腊月初五日 节气:
    calendar_test.go:17: 公历年:2019年12月31日周2 农历年:2019(己亥)[猪]年腊月初六日 节气:
    calendar_test.go:17: 公历年:2020年1月1日周3 农历年:2019(己亥)[猪]年腊月初七日 节气:
    calendar_test.go:17: 公历年:2020年1月2日周4 农历年:2019(己亥)[猪]年腊月初八日 节气:
    calendar_test.go:17: 公历年:2020年1月3日周5 农历年:2019(己亥)[猪]年腊月初九日 节气:
    calendar_test.go:17: 公历年:2020年1月4日周6 农历年:2019(己亥)[猪]年腊月初十日 节气:
    calendar_test.go:17: 公历年:2020年1月5日周0 农历年:2019(己亥)[猪]年腊月十一日 节气:
    calendar_test.go:17: 公历年:2020年1月6日周1 农历年:2019(己亥)[猪]年腊月十二日 节气:小寒
    calendar_test.go:17: 公历年:2020年1月7日周2 农历年:2019(己亥)[猪]年腊月十三日 节气:
    calendar_test.go:17: 公历年:2020年1月8日周3 农历年:2019(己亥)[猪]年腊月十四日 节气:
    calendar_test.go:17: 公历年:2020年1月9日周4 农历年:2019(己亥)[猪]年腊月十五日 节气:
    calendar_test.go:17: 公历年:2020年1月10日周5 农历年:2019(己亥)[猪]年腊月十六日 节气:
    calendar_test.go:17: 公历年:2020年1月11日周6 农历年:2019(己亥)[猪]年腊月十七日 节气:
    calendar_test.go:17: 公历年:2020年1月12日周0 农历年:2019(己亥)[猪]年腊月十八日 节气:
    calendar_test.go:17: 公历年:2020年1月13日周1 农历年:2019(己亥)[猪]年腊月十九日 节气:
    calendar_test.go:17: 公历年:2020年1月14日周2 农历年:2019(己亥)[猪]年腊月廿十日 节气:
    calendar_test.go:17: 公历年:2020年1月15日周3 农历年:2019(己亥)[猪]年腊月廿一日 节气:
    calendar_test.go:17: 公历年:2020年1月16日周4 农历年:2019(己亥)[猪]年腊月廿二日 节气:
    calendar_test.go:17: 公历年:2020年1月17日周5 农历年:2019(己亥)[猪]年腊月廿三日 节气:
    calendar_test.go:17: 公历年:2020年1月18日周6 农历年:2019(己亥)[猪]年腊月廿四日 节气:
    calendar_test.go:17: 公历年:2020年1月19日周0 农历年:2019(己亥)[猪]年腊月廿五日 节气:
    calendar_test.go:17: 公历年:2020年1月20日周1 农历年:2019(己亥)[猪]年腊月廿六日 节气:大寒
    calendar_test.go:17: 公历年:2020年1月21日周2 农历年:2019(己亥)[猪]年腊月廿七日 节气:
    calendar_test.go:17: 公历年:2020年1月22日周3 农历年:2019(己亥)[猪]年腊月廿八日 节气:
    calendar_test.go:17: 公历年:2020年1月23日周4 农历年:2019(己亥)[猪]年腊月廿九日 节气:
    calendar_test.go:17: 公历年:2020年1月24日周5 农历年:2019(己亥)[猪]年腊月卅十日 节气:
    calendar_test.go:17: 公历年:2020年1月25日周6 农历年:2020(庚子)[鼠]年正月初一日 节气:
    calendar_test.go:17: 公历年:2020年1月26日周0 农历年:2020(庚子)[鼠]年正月初二日 节气:
    calendar_test.go:17: 公历年:2020年1月27日周1 农历年:2020(庚子)[鼠]年正月初三日 节气:
    calendar_test.go:17: 公历年:2020年1月28日周2 农历年:2020(庚子)[鼠]年正月初四日 节气:
    calendar_test.go:17: 公历年:2020年1月29日周3 农历年:2020(庚子)[鼠]年正月初五日 节气:
    calendar_test.go:17: 公历年:2020年1月30日周4 农历年:2020(庚子)[鼠]年正月初六日 节气:
    calendar_test.go:17: 公历年:2020年1月31日周5 农历年:2020(庚子)[鼠]年正月初七日 节气:
    calendar_test.go:17: 公历年:2020年2月1日周6 农历年:2020(庚子)[鼠]年正月初八日 节气:
    calendar_test.go:17: 公历年:2020年2月2日周0 农历年:2020(庚子)[鼠]年正月初九日 节气:
    calendar_test.go:17: 公历年:2020年2月3日周1 农历年:2020(庚子)[鼠]年正月初十日 节气:
    calendar_test.go:17: 公历年:2020年2月4日周2 农历年:2020(庚子)[鼠]年正月十一日 节气:立春
    calendar_test.go:17: 公历年:2020年2月5日周3 农历年:2020(庚子)[鼠]年正月十二日 节气:
    calendar_test.go:17: 公历年:2020年2月6日周4 农历年:2020(庚子)[鼠]年正月十三日 节气:
    calendar_test.go:17: 公历年:2020年2月7日周5 农历年:2020(庚子)[鼠]年正月十四日 节气:
    calendar_test.go:17: 公历年:2020年2月8日周6 农历年:2020(庚子)[鼠]年正月十五日 节气:
PASS
```

![](https://github.com/liujiawm/gocalendar/blob/master/test2.png?raw=true)

![](https://github.com/liujiawm/gocalendar/blob/master/test.png?raw=true)
