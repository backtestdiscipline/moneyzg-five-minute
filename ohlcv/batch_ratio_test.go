package ohlcv

import (
	"math"
	"testing"
)

func TestBatchRatioFiveBlocks(t *testing.T) {
	d := make([]Daily, 365)
	for i := range d {
		d[i] = Daily{Net: float64(i/73 - 2), Count: 1}
	}
	mean, ci, e := BatchRatio(d)
	if e != nil || mean != 0 {
		t.Fatal(mean, ci, e)
	}
	se := math.Sqrt(1.25*(73*73*10)) / 365
	if math.Abs(ci.High-2.7764451051977987*se) > 1e-12 || ci.Low != -ci.High {
		t.Fatal(ci)
	}
}
func TestBatchRatioInvalid(t *testing.T) {
	d := make([]Daily, 365)
	for i := range d {
		d[i] = Daily{Net: 1, Count: 1}
	}
	if _, _, e := BatchRatio(d); e == nil {
		t.Fatal("zero variance accepted")
	}
	d[2] = Daily{Net: 1, Count: 0}
	if _, _, e := BatchRatio(d); e == nil {
		t.Fatal("invalid zero count accepted")
	}
	if _, _, e := BatchRatio(d[:364]); e == nil {
		t.Fatal("partial year accepted")
	}
}

func TestBatchRatioIndependentVariableCountReference(t *testing.T) {
	d := make([]Daily, 365)
	for i := range d {
		if i%5 != 0 {
			d[i] = Daily{Net: float64(i%19-9) * 0.13, Count: 1 + i%4}
		}
	}
	theta, ci, e := BatchRatio(d)
	if e != nil {
		t.Fatal(e)
	}
	if math.Abs(theta-(-0.005342465753424658)) > 1e-14 || math.Abs(ci.Low-(-0.012288811615700352)) > 1e-14 || math.Abs(ci.High-0.001603880108851036) > 1e-14 {
		t.Fatal(theta, ci)
	}
}
func TestBatchAdequacy(t *testing.T) {
	d := make([]Daily, 365)
	for i := range d {
		d[i] = Daily{Net: float64(i % 5), Count: 1}
	}
	if ok, e := BatchAdequacy(d, false); !ok || e != nil {
		t.Fatal(ok, e)
	}
	if ok, e := BatchAdequacy(d, true); ok || e != nil {
		t.Fatal("open position accepted")
	}
	for i := 0; i < 73; i++ {
		d[i] = Daily{}
	}
	if ok, e := BatchAdequacy(d, false); ok || e != nil {
		t.Fatal("empty block accepted")
	}
	d[0] = Daily{Net: 1}
	if _, e := BatchAdequacy(d, false); e == nil {
		t.Fatal("invalid day accepted")
	}
}
