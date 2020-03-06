package gocalendar

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCtime(t *testing.T) {
	ct := NewCtime()
	h,m,s := ct.Clock()
	fmt.Printf("%s unix:%d 时:%d : 分:%d : 秒:%d\n",ct.Format(time.RFC3339),ct.UTC().Unix(),h,m,s)
	tn := time.Now()
	h,m,s = tn.Clock()
	fmt.Printf("%s unix:%d 时:%d : 分:%d : 秒:%d\n",tn.UTC().Format(time.RFC3339),tn.Unix(),h,m,s)

	ct = UnixNewCtime(ct.Unix(),0)
	h,m,s = ct.Clock()
	fmt.Printf("%s unix:%d 时:%d : 分:%d : 秒:%d\n",ct.Local().Format(time.RFC3339),ct.UTC().Unix(),h,m,s)
}

func TestDateNewCtime(t *testing.T) {

	// 时区设置，该包放弃使用zoneinfo.zip，默认使用time.Local
	// golang的时区包在$GOROOT/lib/time/zoneinfo.zip
	// 使用时要求设置ENV的ZONEINFO=$GOROOT/lib/time/zoneinfo.zip
	// 如果生产环境中没有，可以设置ENV的ZONEINFO指向该文件实在所在路径
	// 中国时区
	var PRC  = time.FixedZone("CST-8",28800)

	ct := DateNewCtime(&Date{
		Year:2020,
		Month:3,
		Day:5,
		Hour:2,
		Min:46,
		Sec:3,
		Nsec:0,
		Loc:PRC,
	})
	h,m,s := ct.Clock()
	fmt.Printf("%s unix:%d 时:%d : 分:%d : 秒:%d\n",ct.Format(time.RFC3339),ct.UTC().Unix(),h,m,s)

	ct = UnixNewCtime(1583347563,0)
	h,m,s = ct.Clock()
	fmt.Printf("%s unix:%d 时:%d : 分:%d : 秒:%d\n",ct.Format(time.RFC3339),ct.UTC().Unix(),h,m,s)
}

func TestUnixNewCtime(t *testing.T) {
	var u int64 = 1583345522
	ct := UnixNewCtime(u,0)
	fmt.Println(u)
	fmt.Println(ct.Unix())
	fmt.Println(ct.Format("2006 01 02 15:04:05"))
	ln,lo := ct.Zone()
	fmt.Printf("name:%s offset:%d \n",ln,lo)

}

func TestCtime_MonthDay(t *testing.T) {
	t1, _ := time.Parse("2006-01-02 15:04:05", "1583-12-01 00:00:00")
	c := Ctime(t1)
	fmt.Println(c.MonthDay())
}