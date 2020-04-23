package gocalendar

import (
	"fmt"
	"math"
	"time"
)

var (

	// 五行
	WuXingArray = [5]string{"金", "木", "水", "火", "土"}
	// 正五行,天干地支对应五行
	// 天干五行
	TianGanWuXingArray = [10]int{1, 1, 3, 3, 4, 4, 0, 0, 2, 2} // 天干对应五行
	// 地支五行
	DiZhiWuXingArray = [12]int{2, 4, 1, 1, 4, 3, 3, 4, 0, 0, 4, 2} // 地支对应五行

	// 六十花甲子纳音表

	// 十二星座
	ZodiacTitleArray = [12]string{"水瓶", "双鱼", "白羊", "金牛", "双子", "巨蟹", "狮子", "处女", "天秤", "天蝎", "射手", "摩羯"}
	// 星座的起始日期
	ZodiacDayArray = [12]int{20, 19, 21, 20, 21, 22, 23, 23, 23, 24, 22, 22}
)

// isLeap 给定的公历年year是否是闰年
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// LocName 自定义Location名称
func LocName(zone int)string{
	name := "Etc/GMT"
	sign := "-"
	if zone < 0 {
		sign = "+"
		zone = zone * -1
	}

	h,m,_ := time.Unix(int64(zone),0).In(time.UTC).Clock()
	if h == 0 && m == 0 {
		return "UTC"
	}else if m == 0 {
		return fmt.Sprintf("%s%s%d",name,sign,h)
	}
	return fmt.Sprintf("%s%s%d:%d",name,sign,h,m)
}

// TianGanWuXing 天干索引对应五行名称
func TianGanWuXing(i int)string{
	if i <0 || i > 9 {
		return ""
	}

	return WuXingArray[TianGanWuXingArray[i]]
}

// DiZhiWuXing 地支索引对应五行名称
func DiZhiWuXing(i int)string{
	if i <0 || i > 11 {
		return ""
	}

	return WuXingArray[DiZhiWuXingArray[i]]
}

// Zodiac 根据月和日取星座
func Zodiac(m,d int)(int,string){
	i := zodiacIndex(m,d)
	t := zodiacTitle(i)
	return i,t
}
// ZodiacIndex 根据月和日取星座的索引值
func zodiacIndex(m,d int) int {
	k := m - 1

	if d < ZodiacDayArray[k] {
		k = ((k + 12) - 1) % 12
	}

	return k
}
// zodiacTitle 星座名称
func zodiacTitle(i int)string{
	if i < 0 || i > 11 {
		return ""
	}
	return ZodiacTitleArray[i]
}

// B2i bool转int
func B2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
// I2b int转bool
func I2b(i int) bool {
	return i != 0
}

// Pow x的整数n次方
func Pow(x float64, n int) float64 {
	if x == 0 {
		return 0
	}
	result := calPow(x, n)
	if n < 0 {
		result = 1 / result
	}
	return result
}

func calPow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}

	// 向右移动一位
	result := calPow(x, n>>1)
	result *= result

	// 如果n是奇数
	if n&1 == 1 {
		result *= x
	}

	return result
}

// Round 四舍五入保留prec位小数
func Round(n float64, prec int) float64 {
	e := math.Pow10(prec)
	return math.Round(n*e) / e

	// fs := fmt.Sprintf("%."+strconv.Itoa(prec)+"f",n)
	// r,err := strconv.ParseFloat(fs,64)
	// if err != nil {
	// 	return 0
	// }
	// return r
}