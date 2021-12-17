package gocalendar

import (
	"testing"
	"time"
)

func TestJulianDay(t *testing.T) {

	jd := JulianDay(1581,1,1,12)
	if Round(jd,10) == 2298519 {
		t.Log("passed")
	}else{
		t.Error(jd)
	}

	jd = JulianDay(2021,12,6)
	if Round(jd,10) == 2459554.5 {
		t.Log("passed")
	}else{
		t.Error(jd)
	}

	jd = JulianDay(2021,12,6,12)
	if Round(jd,10) == 2459555 {
		t.Log("passed")
	}else{
		t.Error(jd)
	}

	jd = JulianDay(2021,12,6,12,10,10)
	if Round(jd,10) == 2459555.007060185 {
		t.Log("passed")
	}else{
		t.Error(jd)
	}

}

func TestMjd(t *testing.T) {
	jd := Mjd(2021,12,6)
	if Round(jd,10) == 59554 {
		t.Log("passed")
	}else{
		t.Error(jd)
	}

	jd = Mjd(2021,12,6,12)
	if Round(jd,10) == 59554.5 {
		t.Log("passed")
	}else{
		t.Error(jd)
	}

	jd = Mjd(2021,12,6,12,10,10)
	if Round(jd,10) == 59554.5070601851 {
		t.Log("passed")
	}else{
		t.Error(jd)
	}
}

func TestJdToTimeMap(t *testing.T) {
	datetime := JdToTimeMap(2298519)
	if datetime["year"] == 1581 && datetime["month"] == 1 && datetime["day"] == 1 && datetime["hour"] == 12 {
		t.Log("passed")
	}else{
		t.Error(datetime)
	}

	datetime = JdToTimeMap(2459554.5)
	if datetime["year"] == 2021 && datetime["month"] == 12 && datetime["day"] == 6 {
		t.Log("passed")
	}else{
		t.Error(datetime)
	}

	datetime = JdToTimeMap(2459555)
	if datetime["year"] == 2021 && datetime["month"] == 12 && datetime["day"] == 6 && datetime["hour"] == 12 {
		t.Log("passed")
	}else{
		t.Error(datetime)
	}

	datetime = JdToTimeMap(2459555.007060185)
	if datetime["year"] == 2021 && datetime["month"] == 12 && datetime["day"] == 6 && datetime["hour"] == 12 && datetime["minute"] == 10 && datetime["second"] == 10 {
		t.Log("passed")
	}else{
		// fmt.Printf("%d年%d月%d日%d时%d分%d秒\n",datetime["year"], datetime["month"],datetime["day"],datetime["hour"],datetime["minute"],datetime["second"])
		t.Error(datetime)
	}

}

func TestJdToTime(t *testing.T) {
	datetime := JdToTime(2459555.007060185,time.Local)
	t.Log(datetime.Format(time.RFC3339))
}

func TestMjdToTimeMap(t *testing.T) {
	datetime := MjdToTimeMap(59554)
	if datetime["year"] == 2021 && datetime["month"] == 12 && datetime["day"] == 6 {
		t.Log("passed")
	}else{
		t.Error(datetime)
	}

	datetime = MjdToTimeMap(59554.5)
	if datetime["year"] == 2021 && datetime["month"] == 12 && datetime["day"] == 6 && datetime["hour"] == 12 {
		t.Log("passed")
	}else{
		t.Error(datetime)
	}

	datetime = MjdToTimeMap(59554.507060185075)
	if datetime["year"] == 2021 && datetime["month"] == 12 && datetime["day"] == 6 && datetime["hour"] == 12 && datetime["minute"] == 10 && datetime["second"] == 10 {
		t.Log("passed")
	}else{
		t.Error(datetime)
	}
}
