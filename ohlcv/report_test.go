package ohlcv

import "testing"

func TestBootstrapConstantAndDeterministic(t *testing.T) {
	d := make([]Daily, 365)
	for i := range d {
		d[i] = Daily{2, 2}
	}
	a, e := Bootstrap(d)
	if e != nil || a.Low != 1 || a.High != 1 {
		t.Fatalf("%+v %v", a, e)
	}
	b, e := Bootstrap(d)
	if e != nil || a != b {
		t.Fatal("not deterministic")
	}
}
func TestBootstrapNoTradesFails(t *testing.T) {
	if _, e := Bootstrap(make([]Daily, 365)); e == nil {
		t.Fatal("empty resample accepted")
	}
}

// Independent Python SplitMix64 reference for 14 daily observations 0..13,
// seven-day circular blocks and nearest-rank percentiles (10,000 resamples).
func TestBootstrapNonconstantReference(t *testing.T) {
	d := make([]Daily, 14)
	for i := range d {
		d[i] = Daily{Net: float64(i), Count: 1}
	}
	ci, e := Bootstrap(d)
	if e != nil || ci.Low != 4 || ci.High != 9.5 {
		t.Fatalf("interval %+v error %v", ci, e)
	}
}
func TestBootstrapRejectsInconsistentDailyAggregate(t *testing.T) {
	d := make([]Daily, 7)
	d[0] = Daily{Net: 1}
	if _, e := Bootstrap(d); e == nil {
		t.Fatal("nonzero R with no trades accepted")
	}
}
func TestBootstrapBlockSensitivityValidation(t *testing.T) {
	d := make([]Daily, 14)
	for i := range d {
		d[i] = Daily{Net: 1, Count: 1}
	}
	for _, b := range []int{3, 7, 14} {
		ci, e := BootstrapBlocks(d, b)
		if e != nil || ci.Low != 1 || ci.High != 1 {
			t.Fatal(b, ci, e)
		}
	}
	if _, e := BootstrapBlocks(d, 15); e == nil {
		t.Fatal("block exceeds series")
	}
}
