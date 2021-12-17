package gocalendar

import (
	"math"
	"strings"
)



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

// B2i bool转int
//
// true=>1,false=>0
func B2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// I2b int转bool
//
// 0=>false,other=>true
func I2b(i int) bool {
	return i != 0
}

// StringSplice 字符串拼接
func StringSplice(str ...string) string {
	var builder strings.Builder
	var totalLength int

	for _, s := range str {
		totalLength += len(s)
		builder.WriteString(s)
	}

	if totalLength != builder.Len() {
		return ""
	}

	return builder.String()
}