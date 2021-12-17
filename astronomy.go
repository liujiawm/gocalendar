package gocalendar

import (
	"errors"
	"math"
)

const (
	// 均值朔望月长(mean length of Synodic Month)
	cMSM float64 = 29.530588853

	// 以2000年的第一个均值新月点为基准点，此基准点为2000年1月6日14时20分37秒(TT)，其对应真实新月点为2000年1月6日18时13分42秒(TT)
	cBNM float64 = 2451550.0976504628
)

// deltaTDays 地球自转速度调整值Delta T(以∆T表示)
//
// 地球时和UTC的时差 单位:天(days)
func deltaTDays(year,month float64) float64 {
	dt,err := deltaTSeconds(year,month)
	if err != nil {
		return 0
	}

	return Round(dt / 60.0 / 60.0 / 24.0,16)
}

// deltaTMinutes 地球自转速度调整值Delta T(以∆T表示)
//
// 地球时和UTC的时差 单位:分(minutes)
func deltaTMinutes(year,month float64) float64 {
	dt,err := deltaTSeconds(year,month)
	if err != nil {
		return 0
	}

	return Round(dt / 60.0,16)
}

// deltaTSeconds 地球自转速度调整值Delta T(以∆T表示)
//
// 地球时和UTC的时差 单位:秒(seconds)
// 精确至月份
func deltaTSeconds(year,month float64) (float64,error) {
	// 计算方法参考: https://eclipse.gsfc.nasa.gov/SEhelp/deltatpoly2004.html
	// 此算法在-1999年到3000年之间有效

	if year < -1999 || year > 3000 {
		return 0,errors.New("计算DeltaT值限-1999年至3000年之间有效")
	}

	y := year + (month - 0.5) / 12

	var dt float64

	switch {
	case year <= -500:
		u := (year - 1820) / 100
		dt = -20 + 32 * math.Pow(u, 2)
	case year < 500:
		u := y / 100
		dt = 10583.6 - 1014.41*u + 33.78311*math.Pow(u, 2) - 5.952053*math.Pow(u, 3) - 0.1798452*math.Pow(u, 4) + 0.022174192*math.Pow(u, 5) + 0.0090316521*math.Pow(u, 6)
	case year < 1600:
		u := (y - 1000) / 100
		dt = 1574.2 - 556.01*u + 71.23472*math.Pow(u, 2) + 0.319781*math.Pow(u, 3) - 0.8503463*math.Pow(u, 4) - 0.005050998*math.Pow(u, 5) + 0.0083572073*math.Pow(u, 6)
	case year < 1700:
		t := y - 1600
		dt = 120 - 0.9808*t - 0.01532*math.Pow(t, 2) + math.Pow(t, 3)/7129
	case year < 1800:
		t := y - 1700
		dt = 8.83 + 0.1603*t - 0.0059285*math.Pow(t, 2) + 0.00013336*math.Pow(t, 3) - math.Pow(t, 4)/1174000
	case year < 1860:
		t := y - 1800
		dt = 13.72 - 0.332447*t + 0.0068612*math.Pow(t, 2) + 0.0041116*math.Pow(t, 3) - 0.00037436*math.Pow(t, 4) + 0.0000121272*math.Pow(t, 5) - 0.0000001699*math.Pow(t, 6) + 0.000000000875*math.Pow(t, 7)
	case year < 1900:
		t := y - 1860
		dt = 7.62 + 0.5737*t - 0.251754*math.Pow(t, 2) + 0.01680668*math.Pow(t, 3) - 0.0004473624*math.Pow(t, 4) + math.Pow(t, 5)/233174
	case year < 1920:
		t := y - 1900
		dt = -2.79 + 1.494119*t - 0.0598939*math.Pow(t, 2) + 0.0061966*math.Pow(t, 3) - 0.000197*math.Pow(t, 4)
	case year < 1941:
		t := y - 1920
		dt = 21.2 + 0.84493*t - 0.0761*math.Pow(t, 2) + 0.0020936*math.Pow(t, 3)
	case year < 1961:
		t := y - 1950
		dt = 29.07 + 0.407*t - math.Pow(t, 2)/233 + math.Pow(t, 3)/2547
	case year < 1986:
		t := y - 1975
		dt = 45.45 + 1.067*t - math.Pow(t, 2)/260 - math.Pow(t, 3)/718
	case year < 2005:
		t := y - 2000
		dt = 63.86 + 0.3345*t - 0.060374*math.Pow(t, 2) + 0.0017275*math.Pow(t, 3) + 0.000651814*math.Pow(t, 4) + 0.00002373599*math.Pow(t, 5)
	case year < 2050:
		t := y - 2000
		dt = 62.92 + 0.32217*t + 0.005589*math.Pow(t, 2)
	case year < 2150:
		u := (y - 1820) / 100
		dt = -20 + 32*math.Pow(u, 2) - 0.5628*(2150-y)
	default:
		u := (y - 1820) / 100
		dt = -20 + 32*math.Pow(u, 2)
	}

	// 以上的∆T值均假定月球的长期加速度为-26弧秒/cy^2
	// 而Canon中使用的ELP-2000/82月历使用的值略有不同，为-25.858弧秒/cy^2
	// 因此，必须在∆T多项式表达式得出的值上加上一个小的修正“c”，然后才能将其用于标准中
	// 由于1955年至2005年期间的ΔT值是独立于任何月历而得出的，因此该期间无需校正。
	if year < 1955 || year >= 2005 {
		c := -0.000012932 * (y - 1955) * (y - 1955)
		dt += c
	}

	return dt, nil
}


// perturbation 地球在绕日运行时会因受到其他星球之影响而产生摄动(perturbation)
//
// 返回某时刻(儒略日历)的摄动偏移量
func perturbation(jd float64) float64 {
	// 算法公式摘自Jean Meeus在1991年出版的《Astronomical Algorithms》第27章 Equinoxes and solsticesq (第177页)
	// http://www.agopax.it/Libri_astronomia/pdf/Astronomical%20Algorithms.pdf
	// 公式: 0.00001S/∆λ
	// S = Σ[A cos(B+CT)]
	// B和C的单位是度
	// T = JDE0 - J2000 / 36525
	// J2000 = 2451545.0
	// 36525是儒略历一个世纪的天数
	// ∆λ = 1 + 0.0334cosW+0.0007cos2W
	// W = (35999.373T - 2.47)π/180
	// 注释: Liu Min<liujiawm@163.com> https://github.com/liujiawm

	// 公式中A,B,C的值
	ptsA := [24]float64{485, 203, 199, 182, 156, 136, 77, 74, 70, 58, 52, 50, 45, 44, 29, 18, 17, 16, 14, 12, 12, 12, 9, 8}
	ptsB := [24]float64{324.96, 337.23, 342.08, 27.85, 73.14, 171.52, 222.54, 296.72, 243.58, 119.81, 297.17, 21.02, 247.54, 325.15,60.93, 155.12, 288.79, 198.04, 199.76, 95.39, 287.11, 320.81, 227.73, 15.45}
	ptsC := [24]float64{1934.136, 32964.467, 20.186, 445267.112, 45036.886, 22518.443, 65928.934, 3034.906, 9037.513, 33718.147, 150.678, 2281.226, 29929.562, 31555.956, 4443.417, 67555.328, 4562.452, 62894.029, 31436.921, 14577.848, 31931.756, 34777.259, 1222.114, 16859.074}

	T := julianCentury(jd) // T是以儒略世纪(36525日)为单位，以J2000(儒略日2451545.0)为0点

	var s float64 = 0

	for k := 0; k <= 23; k++ {
		s += ptsA[k] * math.Cos(ptsB[k]*2*math.Pi/360+ptsC[k]*2*math.Pi/360*T)
	}

	W := (35999.373*T - 2.47)*2*math.Pi/360

	L := 1 + 0.0334*math.Cos(W) + 0.0007*math.Cos(2*W)

	return Round(0.00001*s/L, 16)
}

// vernalEquinox 计算指定年的春分点
func vernalEquinox(year float64) float64 {
	// 算法公式摘自Jean Meeus在1991年出版的《Astronomical Algorithms》第27章 Equinoxes and solsticesq (第177页)
	// http://www.agopax.it/Libri_astronomia/pdf/Astronomical%20Algorithms.pdf
	// 此公式在-1000年至3000年之间比较准确
	// 在公元前1000年之前或公元3000年之后也可以延申使用，但因外差法求值，年代越远，算出的结果误差就越大。

	var ve float64
	if year >= 1000 && year <= 3000 {
		m := (year - 2000) / 1000
		ve = 2451623.80984 + 365242.37404*m + 0.05169*math.Pow(m, 2) - 0.00411*math.Pow(m, 3) - 0.00057*math.Pow(m, 4)
	}else{
		m := year / 1000
		ve = 1721139.29189 + 365242.1374*m + 0.06134*math.Pow(m, 2) + 0.00111*math.Pow(m, 3) - 0.00071*math.Pow(m, 4)
	}

	return Round(ve, 10)
}

// trueNewMoon 求出实际新月点
//
// 以2000年初的第一个均值新月点为0点求出的均值新月点和其朔望月之序數 k 代入此副程式來求算实际新月点
func trueNewMoon(k float64) float64 {
	// 对于指定的日期时刻JD值jd,算出其为相对于基准点(之后或之前)的第k个朔望月之内。
	// k=INT(jd-bnm)/msm
	// 新月点估值(new moon estimated)为：nme=bnm+msm×k
	// 估计的世纪变数值为：t = (nme - J2000) / 36525
	// 此t是以2000年1月1日12时(TT)为0点，以100年为单位的时间变数，
	// 由于朔望月长每个月都不同，msm所代表的只是其均值，所以算出新月点后，还需要加上一个调整值。
	// adj = 0.0001337×t×t - 0.00000015×t×t×t + 0.00000000073×t×t×t×t
	// 指定日期时刻所属的均值新月点JD值(mean new moon)：mnm=nme+adj

	nme := newMoonEstimated(k)

	t := julianCentury(nme)
	t2 := math.Pow(t, 2) // square for frequent use
	t3 := math.Pow(t, 3) // cube for frequent use
	t4 := math.Pow(t, 4) // to the fourth

	// mean time of phase
	mnm := nme + 0.0001337 * t2 - 0.00000015 * t3 + 0.00000000073 * t4

	// Sun's mean anomaly(地球绕太阳运行均值近点角)(从太阳观察)
	m := 2.5534 + 29.10535669*k - 0.0000218*t2 - 0.00000011*t3

	// Moon's mean anomaly(月球绕地球运行均值近点角)(从地球观察)
	ms := 201.5643 + 385.81693528*k + 0.0107438*t2 + 0.00001239*t3 - 0.000000058*t4

	// Moon's argument of latitude(月球的纬度参数)
	f := 160.7108 + 390.67050274*k - 0.0016341*t2 - 0.00000227*t3 + 0.000000011*t4

	// Longitude of the ascending node of the lunar orbit(月球绕日运行轨道升交点之经度)
	omega := 124.7746 - 1.5637558*k + 0.0020691*t2 + 0.00000215*t3

	// 乘式因子
	e := 1 - 0.002516*t - 0.0000074*t2

	// 因perturbation造成的偏移
	pi180 := math.Pi / 180
	apt1 := -0.4072 * math.Sin(pi180*ms)
	apt1 += 0.17241 * e * math.Sin(pi180*m)
	apt1 += 0.01608 * math.Sin(pi180*2*ms)
	apt1 += 0.01039 * math.Sin(pi180*2*f)
	apt1 += 0.00739 * e * math.Sin(pi180*(ms-m))
	apt1 -= 0.00514 * e * math.Sin(pi180*(ms+m))
	apt1 += 0.00208 * e * e * math.Sin(pi180*(2*m))
	apt1 -= 0.00111 * math.Sin(pi180*(ms-2*f))
	apt1 -= 0.00057 * math.Sin(pi180*(ms+2*f))
	apt1 += 0.00056 * e * math.Sin(pi180*(2*ms+m))
	apt1 -= 0.00042 * math.Sin(pi180*3*ms)
	apt1 += 0.00042 * e * math.Sin(pi180*(m+2*f))
	apt1 += 0.00038 * e * math.Sin(pi180*(m-2*f))
	apt1 -= 0.00024 * e * math.Sin(pi180*(2*ms-m))
	apt1 -= 0.00017 * math.Sin(pi180*omega)
	apt1 -= 0.00007 * math.Sin(pi180*(ms+2*m))
	apt1 += 0.00004 * math.Sin(pi180*(2*ms-2*f))
	apt1 += 0.00004 * math.Sin(pi180*(3*m))
	apt1 += 0.00003 * math.Sin(pi180*(ms+m-2*f))
	apt1 += 0.00003 * math.Sin(pi180*(2*ms+2*f))
	apt1 -= 0.00003 * math.Sin(pi180*(ms+m+2*f))
	apt1 += 0.00003 * math.Sin(pi180*(ms-m+2*f))
	apt1 -= 0.00002 * math.Sin(pi180*(ms-m-2*f))
	apt1 -= 0.00002 * math.Sin(pi180*(3*ms+m))
	apt1 += 0.00002 * math.Sin(pi180*(4*ms))

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

	return Round(mnm + apt1 + apt2,10)
}

// referenceLunarMonthNum 对于指定的日期时刻JD值jd,算出其为相对于基准点(之后或之前)的第几个朔望月
//
// 为从2000年1月6日14时20分36秒起至指定年月日之农历月数,以synodic month为单位
func referenceLunarMonthNum(jd float64) float64{
	return math.Floor((jd - cBNM) / cMSM)
}

// newMoonEstimated 新月点估值(new moon estimated)
func newMoonEstimated(k float64) float64 {
	// 新月点估值(new moon estimated)为：nme=bnm+msm×k
	return cBNM + cMSM * k
}

// meanNewMoon 对于指定日期时刻所属的朔望月,求出其均值新月点的月序数
func meanNewMoon(jd float64) (float64, float64) {
	k := referenceLunarMonthNum(jd)

	nme := newMoonEstimated(k)

	// Time in Julian centuries from 2000 January 0.5.
	t := julianCentury(nme)
	theJd := nme + 0.0001337*math.Pow(t, 2) - 0.00000015*math.Pow(t, 3) + 0.00000000073*math.Pow(t, 4)

	return k, Round(theJd,10)
}

// meanSolarTermsJd 获取指定年以春分开始的24节气(为了确保覆盖完一个公历年，该方法多取2个节气)
//
// 注意：该方法取出的节气时间是未经微调的
func meanSolarTermsJd(year float64) [26]float64 {

	// 该年的春分点jd
	ve := vernalEquinox(year)

	// 该年的回归年长(天)
	// 两个春分点之间为一个回归年长
	ty := vernalEquinox(year + 1) - ve

	ath := 2 * math.Pi / 24

	T := julianThousandYear(ve)
	e := 0.0167086342 - 0.0004203654*T - 0.0000126734*math.Pow(T,2) + 0.0000001444*math.Pow(T,3) - 0.0000000002*math.Pow(T,4) + 0.0000000003*math.Pow(T,5)

	TT := year / 1000
	d := 111.25586939 - 17.0119934518333*TT - 0.044091890166673*math.Pow(TT,2) - 4.37356166661345E-04*math.Pow(TT,3) + 8.16716666602386E-06*math.Pow(TT,4)

	rvp := d * 2 * math.Pi / 360

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

	var mst [26]float64

	for i := 0; i < cap(peri); i++ {
		mst[i] = Round(ve+peri[i]-peri[0], 10)
	}

	return mst
}

// adjustedSolarTermsJd 获取指定年以春分开始的节气
//
// 经过摄动值和deltaT调整后的jd
func adjustedSolarTermsJd(year float64,start, end int) [26]float64 {
	mst := meanSolarTermsJd(year)

	var jqs [26]float64

	for i, jd := range mst {
		if i < start {
			continue
		}
		if i > end {
			continue
		}

		// 取得受perturbation影响所需微调
		pert := perturbation(jd) // perturbation

		// 修正dynamical time to Universal time
		month := math.Floor((float64(i)+1)/2) + 3
		dtd := deltaTDays(year, month) // delta T(天)

		jqs[i] = Round(jd+pert-dtd, 10) // 加上摄动调整值ptb,减去对应的Delta T值(分钟转换为日)
	}

	return jqs
}

// lastYearSolarTerms 取出上一年从冬至开始的6个节气
func lastYearSolarTerms(year float64) [26]float64 {
	return adjustedSolarTermsJd(year-1,18,23)
}


