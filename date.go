// 主要算来来源于www.bieyu.com相关文章
// www.bieyu.com要求最大年限3000年内
// 因此要求年数在1000-3000
package gocalendar

import (
	"fmt"
	"math"
	"time"
)

const PI = math.Pi

type Date struct {
	Year         int
	Month        int
	Day          int
	Hour         int
	Min          int
	Sec          int
	Nsec         int
	Week         int
	Loc          *time.Location
	TianGanDiZhi *TianGanDiZhi // 干支，以立春开始，包含年干支，月干支，日干支，时干支
	JQ           *JQ
}

type LunarDate struct {
	Year      int
	Month     int
	Day       int
	Hour      int
	Min       int
	Sec       int
	Nsec      int
	Week      int
	MonthDays int // 当月有多少天
	LeapYear  int // 是否闰年，0不是闰年，大于就是闰几月
	LeapMonth int // 当前前是否是所闰的那个月，0不是，1本月就是闰月
	Loc       *time.Location
	YearGanZi *LunarYearGanZi // 农历年干支，通俗记年以春节开始
}

type LunarYearGanZi struct {
	Gan     string
	Zhi     string
	Animals string
}

type TianGanDiZhi struct {
	Ytg int // 年天干
	Ydz int // 年地支
	Mtg int // 月天干
	Mdz int // 月地支
	Dtg int // 日天干
	Ddz int // 日地支
	Htg int // 时天干
	Hdz int // 时地支
}

// 节气及时间
type JQ struct {
	Name string
	Date *Date
}

var (
	// 是否区分 早晚子 时,true则23:00-24:00算成上一天
	Zwz = false

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

	// 五行
	WuXingArray = [5]string{"金", "木", "水", "火", "土"}

	// 生肖
	SymbolicAnimalsArray = [12]string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}

	// 十二星座
	XingZuoArray = [12]string{"水瓶", "双鱼", "白羊", "金牛", "双子", "巨蟹", "狮子", "处女", "天秤", "天蝎", "射手", "摩羯"}

	// 星座的起始日期
	XingZuoStartDayArray = [12]int{20, 19, 21, 20, 21, 22, 23, 23, 23, 24, 22, 22}

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

func NewTianGanDiZhi() *TianGanDiZhi {
	return &TianGanDiZhi{}
}

func NewLunarDate(ld *LunarDate) *LunarDate {
	c := NewCtime()
	if ld.Year == 0{
		ld.Year = c.Year()
	}else if ld.Year < 1000 {
		ld.Year = 1000
	}else if ld.Year > 3000 {
		ld.Year = 3000
	}

	if ld.Month < 1 || ld.Month > 12{
		ld.Month = c.Month()
	}

	if ld.Day < 1 {
		ld.Day = c.Day()
	}
	if ld.Day > 29{
		ld.Day = ld.LunarDays()
	}

	if ld.Loc == nil {
		ld.Loc = c.Location()
	}


	return &LunarDate{
		Year:      ld.Year,
		Month:     ld.Month,
		Day:       ld.Day,
		Hour:      ld.Hour,
		Min:       ld.Min,
		Sec:       ld.Sec,
		Nsec:      ld.Nsec,
		Week:      ld.Week,
		LeapYear:  ld.LeapYear,
		LeapMonth: ld.LeapMonth,
		Loc:       ld.Loc,
		YearGanZi: ld.YearGanZi,
	}
}

func NewDate(d *Date) *Date {
	c := NewCtime()
	if d.Year == 0{
		d.Year = c.Year()
	}else if d.Year < 1000 {
		d.Year = 1000
	}else if d.Year > 3000 {
		d.Year = 3000
	}

	if d.Month == 0 {
		d.Month = c.Month()
	}

	if d.Day == 0 {
		d.Day = c.Day()
	}

	if d.Loc == nil {
		d.Loc = c.Location()
	}

	return CtimeNewDate(DateNewCtime(d))
}

func CtimeNewDate(c Ctime) *Date {
	y, m, d := c.Date()
	h, mi, s := c.Clock()

	if y < 1000 {
		y = 1000
	}else if y > 3000 {
		y = 3000
	}

	nd := Date{
		Year:         y,
		Month:        m,
		Day:          d,
		Hour:         h,
		Min:          mi,
		Sec:          s,
		Week:         c.Weekday(),
		Loc:          c.Location(),
		TianGanDiZhi: &TianGanDiZhi{},
	}

	return &nd
}


func (d *Date) Copy() *Date {
	var nd = *d
	return &Date{
		Year:         nd.Year,
		Month:        nd.Month,
		Day:          nd.Day,
		Hour:         nd.Hour,
		Min:          nd.Min,
		Sec:          nd.Sec,
		Nsec:         nd.Nsec,
		Week:         nd.Week,
		Loc:          nd.Loc,
		TianGanDiZhi: nd.TianGanDiZhi,
	}
}

func (d *Date)AddDate(years,months,days int) *Date {
	c := DateNewCtime(d.Copy()).AddDate(years, months, days)

	return CtimeNewDate(c)
}
func (d *Date) PrevYear() *Date {
	return d.AddDate(-1,0,0)
}
func (d *Date) NextYear() *Date {
	return d.AddDate(1,0,0)
}

// 地球在绕日运行时会因受到其他星球之影响而产生摄动(perturbation)
// 返回某时刻(儒略日历)的摄动偏移量
func perturbation(jd float64) float64 {
	t := (jd - 2451545) / 36525
	var s float64 = 0
	for k := 0; k <= 23; k++ {
		s += ptsa[k] * math.Cos(ptsb[k]*2*PI/360+ptsc[k]*2*PI/360*t)
	}
	w := 35999.373*t - 2.47
	l := 1 + 0.0334*math.Cos(w*2*PI/360) + 0.0007*math.Cos(2*w*2*PI/360)
	return round(0.00001*s/l, 16)
}

// 求出实际新月点
// 以2000年初的第一个均值新月点为0点求出的均值新月点和其朔望月之序數 k 代入此副程式來求算实际新月点
func trueNewMoon(k float64) float64 {

	jdt := bnm + k*synMonth

	t := (jdt - 2451545) / 36525 // 2451545为2000年1月1日正午12时的JD

	t2 := t * t // square for frequent use

	t3 := t2 * t // cube for frequent use

	t4 := t3 * t // to the fourth

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
	// 因perturbation造成的偏移：
	apt1 := -0.4072 * math.Sin((PI/180)*mprime)
	apt1 += 0.17241 * es * math.Sin((PI/180)*m)
	apt1 += 0.01608 * math.Sin((PI/180)*2*mprime)
	apt1 += 0.01039 * math.Sin((PI/180)*2*f)
	apt1 += 0.00739 * es * math.Sin((PI/180)*(mprime-m))
	apt1 -= 0.00514 * es * math.Sin((PI/180)*(mprime+m))
	apt1 += 0.00208 * es * es * math.Sin((PI/180)*(2*m))
	apt1 -= 0.00111 * math.Sin((PI/180)*(mprime-2*f))
	apt1 -= 0.00057 * math.Sin((PI/180)*(mprime+2*f))
	apt1 += 0.00056 * es * math.Sin((PI/180)*(2*mprime+m))
	apt1 -= 0.00042 * math.Sin((PI/180)*3*mprime)
	apt1 += 0.00042 * es * math.Sin((PI/180)*(m+2*f))
	apt1 += 0.00038 * es * math.Sin((PI/180)*(m-2*f))
	apt1 -= 0.00024 * es * math.Sin((PI/180)*(2*mprime-m))
	apt1 -= 0.00017 * math.Sin((PI/180)*omega)
	apt1 -= 0.00007 * math.Sin((PI/180)*(mprime+2*m))
	apt1 += 0.00004 * math.Sin((PI/180)*(2*mprime-2*f))
	apt1 += 0.00004 * math.Sin((PI/180)*(3*m))
	apt1 += 0.00003 * math.Sin((PI/180)*(mprime+m-2*f))
	apt1 += 0.00003 * math.Sin((PI/180)*(2*mprime+2*f))
	apt1 -= 0.00003 * math.Sin((PI/180)*(mprime+m+2*f))
	apt1 += 0.00003 * math.Sin((PI/180)*(mprime-m+2*f))
	apt1 -= 0.00002 * math.Sin((PI/180)*(mprime-m-2*f))
	apt1 -= 0.00002 * math.Sin((PI/180)*(3*mprime+m))
	apt1 += 0.00002 * math.Sin((PI/180)*(4*mprime))

	apt2 := 0.000325 * math.Sin((PI/180)*(299.77+0.107408*k-0.009173*t2))
	apt2 += 0.000165 * math.Sin((PI/180)*(251.88+0.016321*k))
	apt2 += 0.000164 * math.Sin((PI/180)*(251.83+26.651886*k))
	apt2 += 0.000126 * math.Sin((PI/180)*(349.42+36.412478*k))
	apt2 += 0.00011 * math.Sin((PI/180)*(84.66+18.206239*k))
	apt2 += 0.000062 * math.Sin((PI/180)*(141.74+53.303771*k))
	apt2 += 0.00006 * math.Sin((PI/180)*(207.14+2.453732*k))
	apt2 += 0.000056 * math.Sin((PI/180)*(154.84+7.30686*k))
	apt2 += 0.000047 * math.Sin((PI/180)*(34.52+27.261239*k))
	apt2 += 0.000042 * math.Sin((PI/180)*(207.19+0.121824*k))
	apt2 += 0.00004 * math.Sin((PI/180)*(291.34+1.844379*k))
	apt2 += 0.000037 * math.Sin((PI/180)*(161.72+24.198154*k))
	apt2 += 0.000035 * math.Sin((PI/180)*(239.56+25.513099*k))
	apt2 += 0.000023 * math.Sin((PI/180)*(331.55+3.592518*k))

	return pt + apt1 + apt2
}

// 对于指定日期时刻所属的朔望月,求出其均值新月点的月序数
func meanNewMoon(jd float64) (float64, float64) {

	// kn为从2000年1月6日14时20分36秒起至指定年月日之阴历月数,以synodic month为单位
	// bnm=2451550.09765为2000年1月6日14时20分36秒之JD值.
	kn := math.Floor((jd - bnm) / synMonth)

	jdt := bnm + kn*synMonth

	// Time in Julian centuries from 2000 January 0.5.
	t := (jdt - 2451545) / 36525
	thejd := jdt + 0.0001337*t*t - 0.00000015*t*t*t + 0.00000000073*t*t*t*t

	return kn, thejd
}

// 计算指定年(公历)的春分点(mean vernal equinox)
// 比利时的气象学家Jean Meeus在1991年出版的”Astronomical Algorithms”一书中提供了一些求均值春分点(mean vernal equinox)的公式
// 但因地球在绕日运行时会因受到其他星球之影响而产生摄动(perturbation),必须将此现象产生的偏移量加入
// 返回儒略日历格林威治时间
func ve(year int) float64 {

	// 这里最大限8000年
	if year < 1000 || year > 8000 {
		fmt.Println("日期超出限制",year)
		return 0
	}

	m := (float64(year) - 2000) / 1000

	s := 2451623.80984 + 365242.37404*m + 0.05169*m*m - 0.00411*m*m*m - 0.00057*m*m*m*m

	return s
	// return round(s,8)
}

// 求∆t
func deltaT(year, month int) float64 {

	y := float64(year) + (float64(month)-0.5)/12

	var dt float64

	if y <= -500 {
		u := (y - 1820) / 100
		dt = -20 + 32*u*u
	} else {
		if y < 500 {
			u := y / 100
			dt = 10583.6 - 1014.41*u + 33.78311*u*u - 5.952053*u*u*u - 0.1798452*u*u*u*u + 0.022174192*u*u*u*u*u + 0.0090316521*u*u*u*u*u*u
		} else {
			if y < 1600 {
				u := (y - 1000) / 100
				dt = 1574.2 - 556.01*u + 71.23472*u*u + 0.319781*u*u*u - 0.8503463*u*u*u*u - 0.005050998*u*u*u*u*u + 0.0083572073*u*u*u*u*u*u
			} else {
				if y < 1700 {
					t := y - 1600
					dt = 120 - 0.9808*t - 0.01532*t*t + t*t*t/7129
				} else {
					if y < 1800 {
						t := y - 1700
						dt = 8.83 + 0.1603*t - 0.0059285*t*t + 0.00013336*t*t*t - t*t*t*t/1174000
					} else {
						if y < 1860 {
							t := y - 1800
							dt = 13.72 - 0.332447*t + 0.0068612*t*t + 0.0041116*t*t*t - 0.00037436*t*t*t*t + 0.0000121272*t*t*t*t*t - 0.0000001699*t*t*t*t*t*t + 0.000000000875*t*t*t*t*t*t*t
						} else {
							if y < 1900 {
								t := y - 1860
								dt = 7.62 + 0.5737*t - 0.251754*t*t + 0.01680668*t*t*t - 0.0004473624*t*t*t*t + t*t*t*t*t/233174
							} else {
								if y < 1920 {
									t := y - 1900
									dt = -2.79 + 1.494119*t - 0.0598939*t*t + 0.0061966*t*t*t - 0.000197*t*t*t*t
								} else {
									if y < 1941 {
										t := y - 1920
										dt = 21.2 + 0.84493*t - 0.0761*t*t + 0.0020936*t*t*t
									} else {
										if y < 1961 {
											t := y - 1950
											dt = 29.07 + 0.407*t - t*t/233 + t*t*t/2547
										} else {
											if y < 1986 {
												t := y - 1975
												dt = 45.45 + 1.067*t - t*t/260 - t*t*t/718
											} else {
												if y < 2005 {
													t := y - 2000
													dt = 63.86 + 0.3345*t - 0.060374*t*t + 0.0017275*t*t*t + 0.000651814*t*t*t*t + 0.00002373599*t*t*t*t*t
												} else {
													if y < 2050 {
														t := y - 2000
														dt = 62.92 + 0.32217*t + 0.005589*t*t
													} else {
														if y < 2150 {
															u := (y - 1820) / 100
															dt = -20 + 32*u*u - 0.5628*(2150-y)
														} else {
															u := (y - 1820) / 100
															dt = -20 + 32*u*u
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if y < 1955 || y >= 2005 {
		dt = dt - (0.000012932 * (y - 1955) * (y - 1955))
	}
	return round(dt/60, 13) // 将秒转换为分
}

// 获取指定年的春分开始的24节气,另外多取2个确保覆盖完一个公历年
// 大致原理是:先用此方法得到理论值,再用摄动值(Perturbation)和固定参数DeltaT做调整
func meanJqJd(year int) [26]float64 {

	// 另外多取2个确保覆盖完一个公历年
	// num := 24 + 2

	jqjd := [26]float64{}

	// 该年的春分點
	jd := ve(year)

	// 该年的回归年长
	ty := ve(year+1) - jd

	ath := 2 * PI / 24

	tx := (jd - 2451545) / 365250

	e := 0.0167086342 - 0.0004203654*tx - 0.0000126734*tx*tx + 0.0000001444*tx*tx*tx - 0.0000000002*tx*tx*tx*tx + 0.0000000003*tx*tx*tx*tx*tx
	tt := float64(year) / 1000

	vp := 111.25586939 - 17.0119934518333*tt - 0.044091890166673*tt*tt - 4.37356166661345E-04*tt*tt*tt + 8.16716666602386E-06*tt*tt*tt*tt
	rvp := vp * 2 * PI / 360

	var peri [26]float64

	for i := 0; i < cap(peri); i++ {
		flag := 0
		th := ath*float64(i) + rvp

		if th > PI && th <= 3*PI {
			th = 2*PI - th
			flag = 1
		}

		if th > 3*PI {
			th = 4*PI - th
			flag = 2
		}

		f1 := 2 * math.Atan(math.Sqrt((1-e)/(1+e))*math.Tan(th/2))

		f2 := (e * math.Sqrt(1-e*e) * math.Sin(th)) / (1 + e*math.Cos(th))

		f := (f1 - f2) * ty / 2 / PI

		if flag == 1 {
			f = ty - f
		}
		if flag == 2 {
			f = 2*ty - f
		}
		peri[i] = f
	}

	for i := 0; i < cap(peri); i++ {
		jqjd[i] = round(jd+peri[i]-peri[0], 8)
	}

	return jqjd
}

// 获取指定年的春分开始作Perturbaton调整后的24节气,可以多取2个
func adjustedJQ(year, start, end int) [26]float64 {

	jq := [26]float64{}

	// 获取该年春分开始的24节气时间点
	jqjd := meanJqJd(year)

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
		jq[k] = round(jd+ptb-dt/60/24+float64(1)/3, 8)
	}

	return jq
}

// 求出以某年立春点开始的节(注意:为了方便计算起运数,此处第0位为上一年的小寒)
func pureJQsinceSpring(year int) [16]float64 {
	jdpjq := [16]float64{}

	dj := adjustedJQ(year-1, 19, 23)

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

	dj = adjustedJQ(year, 0, 25)

	for k, v := range dj {
		if k%2 == 0 || v == 0 {
			continue
		}

		ki++
		jdpjq[ki] = v
	}

	return jdpjq
}

// 求出自冬至点为起点的连续15个中气
func zQsinceWinterSolstice(year int) [15]float64 {
	jdpjq := [15]float64{}
	dj := adjustedJQ(year-1, 18, 23)
	jdpjq[0] = dj[18] // 冬至
	jdpjq[1] = dj[20] // 大寒
	jdpjq[2] = dj[22] // 雨水

	// 求出指定年节气之JD值
	dj = adjustedJQ(year, 0, 23)

	ki := 2 // 数组索引
	for k, v := range dj {
		if k%2 != 0 || v == 0 {
			continue
		}

		ki++
		jdpjq[ki] = v
	}

	return jdpjq
}

// 公历某月有多少天
// 因为golang对1582年10月4日至1582年10月15日加10天
// 这里不再考虑1582年情况
func (d *Date) solarDays() int {

	// 每月多少天的数组，索引为m-1
	md := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	dn := md[d.Month-1]

	if d.Month == 2 && IsLeapYear(d.Year) {
		dn += 1
	}

	return dn
}

// 根据公历月日计算星座下标
func (d *Date)getZodiac() int {
	// 下标从0开始
	kn := d.Month - 1

	if d.Day < XingZuoStartDayArray[kn] {
		kn = ((kn + 12) - 1) % 12
	}

	return kn
}

// 本月1号是星期几
func (d *Date)firstdayWeekday()int{
	nd := d.Copy()
	nd.Day = 1

	return DateNewCtime(nd).Weekday()
}
// 本月最后一天是星期几
func (d *Date)lastdayWeekday()int{
	nd := d.Copy()
	nd.Day = nd.solarDays()

	return DateNewCtime(nd).Weekday()
}
func (d *Date)year() int{
	return d.Year
}

// 星座
func (d *Date)XingZuo() string{
	return XingZuoArray[d.getZodiac()]
}

// 求出含某公历年立春点开始的24节气
// 这里为了日历好显示，取27个节气，从上一年的冬至开始
func (d *Date) Jieqi() map[string]*JQ {

	//jq := [27]JQ{}
	// 为了在日历上显示，将数组改了无序的map
	jq := make(map[string]*JQ,27)

	dj := adjustedJQ(d.PrevYear().year(), 18, 23)

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
		kdate := Julian2Solar(v)
		mapkey := fmt.Sprintf("%d-%d-%d",kdate.Year,kdate.Month,kdate.Day)
		jq[mapkey] = &JQ{
			Name: JieQiArray[jqak],
			Date: kdate,
		}
		/*jq[ji] = JQ{
			Name: JieQiArray[jqak],
			Date: Julian2Solar(v),
		}*/
	}

	dj = adjustedJQ(d.Year, 0, 20)

	for _, v := range dj {

		if v == 0 {
			continue
		}

		ji++
		jqak = (ji + 21) % 24 // 节气名称的索引
		kdate := Julian2Solar(v)
		mapkey := fmt.Sprintf("%d-%d-%d",kdate.Year,kdate.Month,kdate.Day)
		jq[mapkey] = &JQ{
			Name: JieQiArray[jqak],
			Date: kdate,
		}
		/*jq[ji] = JQ{
			Name: JieQiArray[jqak],
			Date: Julian2Solar(v),
		}*/
	}

	return jq
}


// 公历对应的干支
// 以立春时间开始
func (d *Date) GanZhi(){
	GZ := NewTianGanDiZhi()

	dc := d.Copy()

	jd := dc.Solar2Julian()
	if jd == 0 {
		d.TianGanDiZhi = GZ
		return
	}
	jq := pureJQsinceSpring(dc.Year)
	if jd < jq[1] {
		dc = d.PrevYear()
		jq = pureJQsinceSpring(dc.Year)
	}

	tg := [4]int{}
	dz := [4]int{}
	ygz := ((dc.Year+4712+24)%60 + 60) % 60
	GZ.Ytg = ygz % 10 // 年干
	GZ.Ydz = ygz % 12 // 年支

	var ix int
	// 比较求算节气月,求出月干支
	for j := 0; j <= 15; j++ {
		// 已超过指定时刻,故应取前一个节气
		if jq[j] >= jd {
			ix = j - 1
			break
		}
	}

	tmm := ((dc.Year+4712)*12 + (ix - 1) + 60) % 60
	mgz := (tmm + 50) % 60
	GZ.Mtg = mgz % 10 // 月干
	GZ.Mdz = mgz % 12 // 月支

	jda := jd + 0.5                                           // 计算日柱之干支,加0.5是将起始点从正午改为从0点开始
	thes := ((jda - math.Floor(jda)) * 86400) + float64(3600) // 将jd的小数部份化为秒,並加上起始点前移的一小时(3600秒),取其整数值
	dayjd := math.Floor(jda) + thes/86400
	dgz := (int(math.Floor(dayjd+49))%60 + 60) % 60
	GZ.Dtg = dgz % 10 // 日干
	GZ.Ddz = dgz % 12 // 日支
	// 区分早晚子时,日柱前移一柱
	if Zwz && (dc.Hour >= 23) {
		GZ.Dtg = (tg[2] + 10 - 1) % 10
		GZ.Ddz = (dz[2] + 12 - 1) % 12
	}

	dh := dayjd * 12
	hgz := (int(math.Floor(dh+48))%60 + 60) % 60
	GZ.Htg = hgz % 10
	GZ.Hdz = hgz % 12

	d.TianGanDiZhi = GZ

}



// 将公历时间转换为儒略日历时间
func (d *Date) Solar2Julian() float64 {

	yy := float64(d.Year)
	mm := float64(d.Month)
	dd := float64(d.Day)
	hh := float64(d.Hour)
	mi := float64(d.Min)
	ss := float64(d.Sec)

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

	// 1582年10月5日至1582年10月14日这十天是不存在的
	if init == 0 {
		return 0 // 不想加error了,又不知道应该返回什么，先这么着吧
	}

	mp := float64(int(mm+9) % 12)
	jdm := mp*30 + math.Floor((mp+1)*34/57)
	jdd := dd - 1
	jdh := (hh + (mi+(ss/60))/60) / 24
	return round(jdy+jdm+float64(jdd)+float64(jdh)+init, 7)

}

// 将儒略日历时间转换为公历(格里高利历)时间
func Julian2Solar(jd float64) *Date {
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
	sd := math.Floor(round((jd+0.5-math.Floor(jd+0.5))*24, 5)*60*60 + 0.00005)

	mt := math.Floor(sd / 60)
	ss := int(sd) % 60
	hh := math.Floor(mt / 60)
	mt = float64(int(mt) % 60)
	yy := math.Floor(y)
	mm := math.Floor(m)
	dd := math.Floor(da)

	return  NewDate(&Date{
		Year:  int(yy),
		Month: int(mm),
		Day:   int(dd),
		Hour:  int(hh),
		Min:   int(mt),
		Sec:   int(ss),
	})
}

// 以比较日期法求算冬月及其余各月名称代码,包含闰月,冬月为0,腊月为1,正月为2,其余类推.闰月多加0.5
func zQandSMandLunarMonthCode(year int) ([15]float64, [16]float64, [15]float64) {

	mc := [15]float64{}
	// 求出自冬至点为起点的连续15个中气
	jdzq := zQsinceWinterSolstice(year)

	// 求出以含冬至中气为阴历11月(冬月)开始的连续16个朔望月的新月點
	jdnm := sMsinceWinterSolstice(year, jdzq[0])

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

	return jdzq, jdnm, mc
}

// 求算以含冬至中气为阴历11月开始的连续16个朔望月
func sMsinceWinterSolstice(year int, jdws float64) [16]float64 {

	tjd := [20]float64{}
	jdnm := [16]float64{}

	// 求年初前两个月附近的新月点(即前一年的11月初)
	d := NewDate(&Date{
		Year:  year,
		Month: 11,
		Day:   1,
		Hour:  0,
		Min:   0,
		Sec:   0,
	})
	jd := d.PrevYear().Solar2Julian()

	// 求得自2000年1月起第kn个平均朔望日及其JD值
	// kn,thejd := meanNewMoon(jd)
	kn, _ := meanNewMoon(jd)

	// 求出连续20个朔望月
	for i := 0; i <= 19; i++ {
		k := kn + float64(i)

		// mjd := thejd + synMonth * float64(i)

		// 以k值代入求瞬時朔望日,因中國比格林威治先行8小時,加1/3天
		tjd[i] = trueNewMoon(k) + float64(1)/3

		// 下式为修正dynamical time to Universal time
		month := i - 1
		tjd[i] = round(tjd[i]-deltaT(year, month)/1440, 7) // 1为1月,0为前一年12月,-1为前一年11月(当i=0时,i-1=-1,代表前一年11月)
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
	return jdnm

}

// 获取农历某年的闰月,0为无闰月
func (ld *LunarDate) leap() int {

	_, _, mc := zQandSMandLunarMonthCode(ld.Year)

	var leap float64 = 0 // 若闰月旗标为0代表无闰月
	for j := 1; j <= 14; j++ { // 确认指定年前一年11月开始各月是否闰月
		if mc[j]-math.Floor(mc[j]) > 0 { // 若是,则将此闰月代码放入闰月旗标內
			leap = math.Floor(mc[j] + 0.5) // leap = 0对应农历11月,1对应农历12月,2对应农历隔年1月,依此类推.
			break
		}
	}

	lm := int(math.Max(0, leap-2))
	ld.LeapYear = lm // 是否闰年，如果是闰年，这里就是所闰的月

	return int(math.Max(0, leap-2))
}

// 获取农历某个月有多少天
func (ld *LunarDate) LunarDays() int {

	_, jdnm, mc := zQandSMandLunarMonthCode(ld.Year)

	leap := 0 // 若闰月旗标为0代表无闰月
	for j := 1; j <= 14; j++ { // 确认指定年前一年11月开始各月是否闰月
		if mc[j]-math.Floor(mc[j]) > 0 { // 若是,则将此闰月代码放入闰月旗标內
			leap = int(math.Floor(mc[j] + 0.5)) // leap=0对应农历11月,1对应农历12月,2对应农历隔年1月,依此类推.
			break
		}
	}

	// 11月对应到1,12月对应到2,1月对应到3,2月对应到4,依此类推
	mm := ld.Month + 2 // 用农历月份简单推出公历月份

	// 求算农历各月之大小,大月30天,小月29天
	nofd := [15]int{}
	for i := 0; i <= 14; i++ {
		nofd[i] = int(math.Floor(jdnm[i+1]+0.5) - math.Floor(jdnm[i]+0.5)) // 每月天数,加0.5是因JD以正午起算
	}

	dy := 0 // 当月天数

	if ld.LeapMonth == 1 { // 若是闰月
		if leap < 3 { // 而旗标非闰月或非本年闰月,则表示此年不含闰月.leap=0代表无闰月,=1代表闰月为前一年的11月,=2代表闰月为前一年的12月
			return 0
		} else { // 若本年內有闰月
			if leap != mm { // 但不为指定的月份
				return 0
			} else { // 若指定的月份即为闰月
				dy = nofd[mm]
			}
		}
	} else { // 若沒有指明是闰月
		if leap == 0 { // 若旗标非闰月,则表示此年不含闰月(包括前一年的11月起之月份)
			dy = nofd[mm-1]
		} else { // 若旗标为本年有闰月(包括前一年的11月起之月份) 公式nofd(mx - (mx > leap) - 1)的用意为:若指定月大于闰月,则索引用mx,否则索引用mx-1
			dy = nofd[mm+B2i(mm > leap)-1]
		}
	}

	return dy
}

// 农历年干支通俗记法，以春节开始
func (ld *LunarDate)LunarYearGanZiCommon(){
	gk := (ld.Year - 3) % 10
	zk := (ld.Year - 3) % 12

	if gk == 0 {
		gk = 9
	}else {
		gk -= 1
	}

	if zk == 0 {
		zk = 11
	}else {
		zk -= 1
	}
	ld.YearGanZi = &LunarYearGanZi{
		Gan:TianGanArray[gk],
		Zhi:DiZhiArray[zk],
		Animals:SymbolicAnimalsArray[zk],
	}
}

// 将农历时间转换成公历时间
func (ld *LunarDate) Lunar2Solar() *Date {

	_, jdnm, mc := zQandSMandLunarMonthCode(ld.Year)

	leap := 0 // 若闰月旗标为0代表无闰月
	for j := 1; j <= 14; j++ { // 确认指定年前一年11月开始各月是否闰月
		if mc[j]-math.Floor(mc[j]) > 0 { // 若是,则将此闰月代码放入闰月旗标內
			leap = int(math.Floor(mc[j] + 0.5)) // leap=0对应农历11月,1对应农历12月,2对应农历隔年1月,依此类推.
			break
		}
	}

	// 11月对应到1,12月对应到2,1月对应到3,2月对应到4,依此类推
	mm := ld.Month + 2 // 用农历月份简单推出公历月份

	// 求算农历各月之大小,大月30天,小月29天
	nofd := [15]int{}
	for i := 0; i <= 14; i++ {
		nofd[i] = int(math.Floor(jdnm[i+1]+0.5) - math.Floor(jdnm[i]+0.5)) // 每月天数,加0.5是因JD以正午起算
	}

	var jd float64 = 0

	if ld.LeapMonth == 1 { // 若是闰月
		if leap < 3 { // 而旗标非闰月或非本年闰月,则表示此年不含闰月.leap=0代表无闰月,=1代表闰月为前一年的11月,=2代表闰月为前一年的12月
			return new(Date) // 此年非闰年
		} else { // 若本年內有闰月
			if leap != mm { // 但不为指入的月份
				return new(Date) // 则指定的月份非闰月,此月非闰月
			} else { // 若输入的月份即为闰月
				if ld.Month <= nofd[mm] { // 若指定的日期不大于当月的天數
					jd = jdnm[mm] + float64(ld.Day) - 1 // 则将当月之前的JD值加上日期之前的天數
				} else { // 日期超出范围
					return new(Date)
				}
			}
		}
	} else { // 若沒有勾選閏月則
		if leap == 0 { // 若旗標非閏月,則表示此年不含閏月(包括前一年的11月起之月份)
			if ld.Day <= nofd[mm-1] { // 若輸入的日期不大於當月的天數
				jd = jdnm[mm-1] + float64(ld.Day) - 1 // 則將當月之前的JD值加上日期之前的天數
			} else { // 日期超出範圍
				return new(Date)
			}
		} else { // 若旗標為本年有閏月(包括前一年的11月起之月份) 公式nofd(mx - (mx > leap) - 1)的用意為:若指定月大於閏月,則索引用mx,否則索引用mx-1
			if ld.Day <= nofd[mm+B2i(mm > leap)-1] { // 若輸入的日期不大於當月的天數
				jd = jdnm[mm+B2i(mm > leap)-1] + float64(ld.Day) - 1 // 則將當月之前的JD值加上日期之前的天數
			} else { // 日期超出範圍
				return new(Date)
			}
		}
	}

	// 因为儒略日历时间需要用公历计算，所以取回公历的年月日之后，直接把ld的时间替换回去
	dT := Julian2Solar(jd)
	dT.Hour = ld.Hour
	dT.Min = ld.Min
	dT.Sec = ld.Sec
	dT.Loc = ld.Loc

	dT.GanZhi() // 重新取干支(因为时间有变化)

	return dT
}

// 将公历时间转换成农历时间
func (d *Date) Solar2Lunar() *LunarDate {
	isLeapMonth := 0
	// 求出指定年月日之JD值
	dTemp := d.Copy()
	dTemp.Hour = 12
	dTemp.Min = 0
	dTemp.Sec = 0
	jd := dTemp.Solar2Julian()

	prev := 0
	_, jdnm, mc := zQandSMandLunarMonthCode(d.Year)
	if math.Floor(jd) < math.Floor(jdnm[0]+0.5) {
		prev = 1
		_, jdnm, mc = zQandSMandLunarMonthCode(d.Year - 1)
	}

	var mi = 0
	for i := 0; i <= 14; i++ { // 指令中加0.5是为了改为从0时算起而不从正午算起
		if math.Floor(jd) >= math.Floor(jdnm[i]+0.5) && math.Floor(jd) < math.Floor(jdnm[i+1]+0.5) {
			mi = i
			break
		}
	}

	ld := NewLunarDate(&LunarDate{
		Year:  d.Year,
		Month: d.Month,
		Day:   d.Day,
		Hour:  d.Hour,
		Min:   d.Min,
		Sec:   d.Sec,
		Loc:   d.Loc,
	})

	if mc[mi] < 2 || prev == 1 { // 年
		ld.Year = d.Year - 1
	}

	lm := ld.leap() // 闰几月，0为无闰月

	if lm > 0 && (mc[mi]-math.Floor(mc[mi]))*2+1 != 1 { // 因mc(mi)=0对应到前一年农历11月,mc(mi)=1对应到前一年农历12月,mc(mi)=2对应到本年1月,依此类推
		isLeapMonth = 1
	}

	ld.Month = int(math.Floor(mc[mi]+10))%12 + 1 // 月

	if isLeapMonth == 1 {
		ld.LeapMonth = 1 // 当前月是闰月
	}

	ld.MonthDays = ld.LunarDays()

	ld.Day = int(math.Floor(jd) - math.Floor(jdnm[mi]+0.5) + 1) // 日,此处加1是因为每月初一从1开始而非从0开始

	ld.LunarYearGanZiCommon()
	return ld
}

// 农历月份常用名称
func MonthChinese(m int)string{
	if m > 0 && m <= 12 {
		return MonthChineseArray[m-1]
	}
	return ""
}

// 农历日期数字返回汉字表示法
func DayChinese(d int)string{
	daystr := ""
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



// 四舍五入保留s位小数
func round(n float64, prec int) float64 {
	e := math.Pow10(prec)
	return math.Round(n*e) / e

	/*fs := fmt.Sprintf("%."+strconv.Itoa(s)+"f",n)
	r,e := strconv.ParseFloat(fs,64)
	if e != nil {
		return 0
	}
	return r*/
}
