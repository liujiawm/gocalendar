package gocalendar

import (
	"testing"
	"time"
)

func TestIsLeapYear(t *testing.T) {
	if !IsLeapYear(1900) && IsLeapYear(1940) && IsLeapYear(2000) && !IsLeapYear(2100) {
		t.Log("passed")
	}else{
		t.Error()
	}
}

func TestBeginningOfMonth(t *testing.T) {
	ti := time.Date(2021,5,15,0,0,0,10,time.UTC)
	tb := time.Date(2021,5,1,0,0,0,10,time.UTC)
	if BeginningOfMonth(ti).Equal(tb){
		t.Log("passed")
	}else{
		t.Error()
	}
}

func TestEndOfMonth(t *testing.T) {
	ti := time.Date(2021,5,15,0,0,0,10,time.UTC)
	tb := time.Date(2021,5,31,0,0,0,10,time.UTC)
	if EndOfMonth(ti).Equal(tb){
		t.Log("passed")
	}else{
		t.Error()
	}
}
