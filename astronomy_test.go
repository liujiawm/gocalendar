package gocalendar

import (
	"testing"
)


func TestDeltaTDays(t *testing.T) {
	dtd := deltaTDays(2021,12)
	// 0.0008406386097956
	if Round(dtd,10) == 0.0008406386 {
		t.Log("passed")
	}else{
		t.Error(dtd)
	}

}


func TestDeltaTMinutes(t *testing.T) {
	dtmi := deltaTMinutes(2021,12)
	// 1.2105195981056707
	if Round(dtmi,10) == 1.2105195981 {
		t.Log("passed")
	}else{
		t.Error(dtmi)
	}

}


func TestPerturbation(t *testing.T){
	var jd float64 = 2298519
	pt := perturbation(jd)
	// -0.0056141748866219
	if Round(pt,10) == -0.0056141749 {
		t.Log("passed")
	}else{
		t.Error(pt)
	}
}

func TestVernalEquinox(t *testing.T){
	ve := vernalEquinox(2021)
	if Round(ve,10) == 2459293.8997175973 {
		t.Log("passed")
	}else{
		t.Error(ve)
	}
}

func TestMeanSolarTerms(t *testing.T){
	mst := meanSolarTermsJd(2021)
	if Round(mst[0],10) == 2459293.8997175973 && Round(mst[1],10) == 2459309.060708575 && mst[2] == 2459324.354092278 && Round(mst[24],10) == 2459659.142093854 && Round(mst[25],10) == 2459674.3030848317 {
		t.Log("passed")
	}else{
		t.Error(mst)
	}
}

func TestTrueNewMoon(t *testing.T){
	var jd float64 = 2298519
	k := referenceLunarMonthNum(jd)
	tnm := trueNewMoon(k)
	if Round(tnm,10) == 2298493.2989711817 {
		t.Log("passed")
	}else{
		t.Error(tnm)
	}
}

func TestAdjustedSolarTermsJd(t *testing.T){
	jqs := adjustedSolarTermsJd(2021,0,25)

	if Round(jqs[0],10) == 2459293.9010286564 && Round(jqs[1],10) == 2459309.0658356417 && jqs[2] == 2459324.356054907 && Round(jqs[24],10) == 2459659.1481248834 && Round(jqs[25],10) == 2459674.3054912435 {
		t.Log("passed")
	}else{
		t.Error(jqs)
	}
}

func TestLastYearSolarTerms(t *testing.T){
	ljqs := lastYearSolarTerms(2021)
	if ljqs[0] == 0 && ljqs[17] == 0 && ljqs[24] == 0 && Round(ljqs[18],10) == 2459204.9184778044 && Round(ljqs[23],10) == 2459278.8707997804 {
		t.Log("passed")
	}else{
		t.Error(ljqs)
	}
}

