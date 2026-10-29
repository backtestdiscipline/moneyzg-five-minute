package ohlcv

import (
	"math"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
func TestStrictFractalTimingAndEMAConditions(t *testing.T) {
	bars := make([]Bar, Warmup+6)
	e20 := make([]float64, len(bars))
	e50 := make([]float64, len(bars))
	for i := range bars {
		bars[i] = Bar{Low: 13, High: 15}
		e20[i] = 13
		e50[i] = 10
	}
	p := Warmup + 2
	bars[p].Low = 12
	if makeSignal(bars, e20, e50, p+1).dir != 0 {
		t.Fatal("emitted before second right bar")
	}
	if makeSignal(bars, e20, e50, p+2).dir != 1 {
		t.Fatal("missing confirmed long")
	}
	bars[p+2].Low = 12
	if makeSignal(bars, e20, e50, p+2).dir != 0 {
		t.Fatal("tie allowed")
	}
	bars[p+2].Low = 13
	e20[p+2] = 9
	if makeSignal(bars, e20, e50, p+2).dir != 0 {
		t.Fatal("confirmation trend ignored")
	}
	e20[p+2] = 13
	e50[p] = 12
	if makeSignal(bars, e20, e50, p+2).dir != 0 {
		t.Fatal("EMA touch accepted")
	}
	for i := range bars {
		bars[i].Low = 9
		bars[i].High = 11
		e20[i] = 10
		e50[i] = 14
	}
	bars[p].High = 12
	if makeSignal(bars, e20, e50, p+2).dir != -1 {
		t.Fatal("missing confirmed short")
	}
}
func TestEntryExitGapConflictAndCosts(t *testing.T) {
	b := Bar{Time: 1, Open: 100, High: 110, Low: 90}
	p := enter(b, signal{1, 95})
	if p == nil || !near(p.trade.Entry, 100.02) || !near(p.trade.Target, 107.55) {
		t.Fatalf("entry %+v", p)
	}
	raw, reason, amb, ok := exit(b, p.trade)
	if !ok || raw != 95 || reason != "stop" || !amb {
		t.Fatalf("ambiguous exit %v %v %v", raw, reason, amb)
	}
	finish(&p.trade, p.rawEntry, raw, 0.1, 2, reason, amb)
	if !near(p.trade.Exit, 94.981) || !near(p.trade.Fees, (100.02+94.981)*Fee) || !near(p.trade.NetR, (-5.039-p.trade.Fees-0.1)/5.02) {
		t.Fatalf("accounting %+v", p.trade)
	}
	b.Open = 90
	raw, _, _, _ = exit(b, p.trade)
	if raw != 90 {
		t.Fatal("long gap stop should fill at worse open")
	}
	s := enter(Bar{Open: 100}, signal{-1, 105})
	raw, why, _, ok := exit(Bar{Open: 108, High: 110, Low: 107}, s.trade)
	if !ok || raw != 108 || why != "stop" {
		t.Fatal("short gap")
	}
	raw, why, _, ok = exit(Bar{Open: 100, High: 100, Low: 90}, s.trade)
	if !ok || !near(raw, s.trade.Target) || why != "target" {
		t.Fatal("short target")
	}
	if enter(Bar{Open: 100}, signal{1, 101}) != nil || enter(Bar{Open: 100}, signal{-1, 99}) != nil {
		t.Fatal("invalid risk accepted")
	}
}
func TestDataAndFundingFailClosed(t *testing.T) {
	start := int64(86400000)
	bars := make([]Bar, Warmup+3)
	for i := range bars {
		bars[i] = Bar{Time: start - int64(Warmup-i)*Interval, Open: 100, High: 101, Low: 99, Close: 100, Volume: 1}
	}
	end := start + 3*Interval
	if _, err := Evaluate(bars, nil, start, end); err == nil {
		t.Fatal("missing settlement accepted")
	}
	funding := []Funding{{Time: start, Rate: 0.001, Mark: 100}}
	if _, err := Evaluate(bars, funding, start, end); err != nil {
		t.Fatal(err)
	}
	bars[1].Time++
	if _, err := Evaluate(bars, funding, start, end); err == nil {
		t.Fatal("gap accepted")
	}
}
func TestPositionAccountingFundingAndDirection(t *testing.T) {
	for _, dir := range []int{1, -1} {
		stop := 95.0
		if dir == -1 {
			stop = 105
		}
		p := enter(Bar{Open: 100}, signal{dir, stop})
		fund := float64(dir) * 0.001 * 100
		finish(&p.trade, p.rawEntry, 100, fund, 1, "test", false)
		if !near(p.trade.FundingCost, fund) || p.trade.Fees <= 0 || p.trade.SlippageCost <= 0 {
			t.Fatalf("costs %+v", p.trade)
		}
	}
}

func TestEndToEndEntryOverlapAndNoFutureLeak(t *testing.T) {
	start := int64(86400000 + Interval)
	n := Warmup + 16
	bars := make([]Bar, n)
	for i := range bars {
		c := 100 + float64(i)*0.1
		bars[i] = Bar{Time: start + int64(i-Warmup)*Interval, Open: c, High: c + 0.1, Low: c - 0.1, Close: c, Volume: 1}
	}
	p := Warmup + 2
	bars[p].Low = bars[p].Close - 1.5
	bars[p+4].Low = bars[p+4].Close - 1.5
	bars[p+9].Low = bars[p+9].Close - 4
	end := start + int64(n-Warmup)*Interval
	result, err := Evaluate(bars, nil, start, end)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Trades) != 1 || result.Trades[0].EntryTime != bars[p+3].Time || result.IgnoredSignals < 1 || result.OpenPosition {
		t.Fatalf("unexpected result %+v", result)
	}
	altered := append([]Bar(nil), bars...)
	altered[p+11].High += 20
	again, err := Evaluate(altered, nil, start, end)
	if err != nil {
		t.Fatal(err)
	}
	if len(again.Trades) < 1 || result.Trades[0] != again.Trades[0] {
		t.Fatal("future data changed earlier trade")
	}
	shortEnd := start + int64(p+4-Warmup)*Interval
	partial, err := Evaluate(bars[:p+4], nil, start, shortEnd)
	if err != nil {
		t.Fatal(err)
	}
	if !partial.OpenPosition || len(partial.Trades) != 0 {
		t.Fatalf("boundary not incomplete: %+v", partial)
	}
}

func TestOffGridFundingTimestampFailsClosed(t *testing.T) {
	start := int64(86400000)
	bars := make([]Bar, Warmup+1)
	for i := range bars {
		bars[i] = Bar{Time: start + int64(i-Warmup)*Interval, Open: 100, High: 101, Low: 99, Close: 100, Volume: 1}
	}
	if _, err := Evaluate(bars, []Funding{{Time: start + 1, Rate: 0.001, Mark: 100}}, start, start+Interval); err == nil {
		t.Fatal("off-grid funding timestamp must not be silently rounded")
	}
}
