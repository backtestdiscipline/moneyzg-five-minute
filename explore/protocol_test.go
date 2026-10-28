package explore

import "testing"

// These explicit examples are protocol review fixtures, not the eventual evaluator.
func TestStrictPivotAndTrendConventions(t *testing.T) {
	type fixture struct {
		name                                               string
		lows                                               []int
		pivotEMA20, pivotEMA50, confirmEMA20, confirmEMA50 int
		want                                               bool
	}
	cases := []fixture{
		{"valid", []int{14, 13, 11, 13, 14}, 12, 10, 13, 11, true},
		{"equal neighbor", []int{14, 13, 11, 11, 14}, 12, 10, 13, 11, false},
		{"touches EMA", []int{14, 13, 10, 13, 14}, 12, 10, 13, 11, false},
		{"trend reversed at confirmation", []int{14, 13, 11, 13, 14}, 12, 10, 10, 11, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pivots := confirmed(c.lows, 4, false)
			got := len(pivots) == 1 && c.pivotEMA50 < c.lows[2] && c.lows[2] < c.pivotEMA20 && c.confirmEMA20 > c.confirmEMA50
			if got != c.want {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestEntryDelayAndObservedQuoteExit(t *testing.T) {
	// Pivot starts at minute 0; confirmation is close of [10,15), entry
	// earliest after 15:01. Quotes are observations, not generated fills.
	type quote struct {
		second   int
		bid, ask int
	}
	quotes := []quote{{900, 100, 101}, {901, 101, 102}, {930, 98, 99}, {931, 97, 98}}
	eligible := 15*60 + 1
	entry := quote{}
	for _, q := range quotes {
		if q.second >= eligible {
			entry = q
			break
		}
	}
	if entry.second != 901 || entry.ask != 102 {
		t.Fatalf("wrong eligible entry: %+v", entry)
	}
	// Fixed stop 99; first bid reaches it at 930, but fill is first observed
	// bid at or after 931, reflecting the specified one-second order delay.
	trigger, fill := quote{}, quote{}
	for _, q := range quotes {
		if q.second > entry.second && q.bid <= 99 {
			trigger = q
			break
		}
	}
	for _, q := range quotes {
		if q.second >= trigger.second+1 {
			fill = q
			break
		}
	}
	if trigger.second != 930 || fill.second != 931 || fill.bid != 97 {
		t.Fatalf("trigger %+v fill %+v", trigger, fill)
	}
	if entry.ask <= 99 || 102.0+1.5*(102.0-99.0) != 106.5 {
		t.Fatal("invalid risk fixture")
	}
}
