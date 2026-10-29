// Package ohlcv implements a study-specific, exploratory candle simulation.
// It does not claim that bar prices are executable exchange fills.
package ohlcv

import (
	"fmt"
	"math"
)

const (
	Interval int64 = 300000 // milliseconds
	Warmup         = 500
	Fee            = 0.0005
	Slippage       = 0.0002
)

type Bar struct {
	Time                           int64
	Open, High, Low, Close, Volume float64
}

// Funding retains the source timestamp in milliseconds.
type Funding struct {
	Time       int64
	Rate, Mark float64
}
type Trade struct {
	Direction                                     int
	EntryTime, ExitTime                           int64
	Entry, Exit, Stop, Target, Risk               float64
	GrossR, NetR, Fees, SlippageCost, FundingCost float64
	Ambiguous, FundingAmbiguous                   bool
	ExitReason                                    string
}
type Result struct {
	Trades                                          []Trade
	InvalidEntries, IgnoredSignals, AmbiguousTrades int
	OpenPosition                                    bool
}
type signal struct {
	dir  int
	stop float64
}
type position struct {
	trade             Trade
	rawEntry, funding float64
}

func finite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }
func validBar(b Bar) bool {
	return b.Time%Interval == 0 && finite(b.Open) && finite(b.High) && finite(b.Low) && finite(b.Close) && finite(b.Volume) && b.Low > 0 && b.Volume >= 0 && b.High >= b.Low && b.Open >= b.Low && b.Open <= b.High && b.Close >= b.Low && b.Close <= b.High
}

// Evaluate takes precisely 500 warm-up bars followed by the evaluation partition.
// start is the first evaluation bar; end is exclusive. Funding must contain exactly
// every 00:00, 08:00 and 16:00 UTC boundary in [start,end). Missing data fail closed.
func Evaluate(bars []Bar, funding []Funding, start, end int64) (Result, error) {
	var out Result
	if start%Interval != 0 || end%Interval != 0 || start >= end || len(bars) < Warmup || int64(len(bars)-Warmup) != (end-start)/Interval {
		return out, fmt.Errorf("invalid partition or warm-up length")
	}
	for i, b := range bars {
		if !validBar(b) || b.Time != start-int64(Warmup-i)*Interval {
			return out, fmt.Errorf("invalid or noncontiguous bar %d", i)
		}
	}
	const settlement int64 = 8 * 60 * 60 * 1000
	funds := make(map[int64]Funding)
	for _, f := range funding {
		if f.Time < start || f.Time >= end || f.Time%settlement > 1000 || !finite(f.Rate) || !finite(f.Mark) || f.Mark <= 0 {
			return out, fmt.Errorf("invalid funding")
		}
		nominal := f.Time - f.Time%settlement
		if _, ok := funds[nominal]; ok {
			return out, fmt.Errorf("duplicate funding")
		}
		funds[nominal] = f
	}
	first := start
	if r := first % settlement; r != 0 {
		first += settlement - r
	}
	for t := first; t < end; t += settlement {
		if _, ok := funds[t]; !ok {
			return out, fmt.Errorf("missing funding at %d", t)
		}
	}
	e20 := make([]float64, len(bars))
	e50 := make([]float64, len(bars))
	e20[0] = bars[0].Close
	e50[0] = bars[0].Close
	var pos *position
	var pending *signal
	for i, b := range bars {
		if i > 0 {
			e20[i] = e20[i-1] + 2.0/21*(b.Close-e20[i-1])
			e50[i] = e50[i-1] + 2.0/51*(b.Close-e50[i-1])
		}
		if i < Warmup {
			continue
		}
		f, hasFunding := funds[b.Time]
		// Exact-timestamp ties use existing exposure before open exits or entries.
		if pos != nil && hasFunding && f.Time == b.Time {
			pos.funding += float64(pos.trade.Direction) * f.Rate * f.Mark
		}
		// A position already open can exit at the open before a later funding event.
		if pos != nil {
			if raw, reason, amb, ok := openExit(b, pos.trade); ok {
				finish(&pos.trade, pos.rawEntry, raw, pos.funding, b.Time, reason, amb)
				out.Trades = append(out.Trades, pos.trade)
				pos = nil
			}
		}
		if pending != nil {
			if pos != nil {
				panic("pending with open position")
			}
			pos = enter(b, *pending)
			if pos == nil {
				out.InvalidEntries++
			}
			pending = nil
		}
		if pos != nil {
			if raw, reason, amb, ok := openExit(b, pos.trade); ok {
				finish(&pos.trade, pos.rawEntry, raw, pos.funding, b.Time, reason, amb)
				out.Trades = append(out.Trades, pos.trade)
				pos = nil
			}
		}
		if pos != nil {
			raw, reason, amb, ok := exit(b, pos.trade)
			if hasFunding && f.Time > b.Time {
				charge := float64(pos.trade.Direction) * f.Rate * f.Mark
				// If an intrabar exit and settlement coexist, order is unknown.
				// Select the adverse funding outcome, and disclose ambiguity.
				cost, ambiguous := fundingOutcome(charge, ok)
				applyFunding(pos, cost, ambiguous)
			}
			if ok {
				finish(&pos.trade, pos.rawEntry, raw, pos.funding, b.Time, reason, amb)
				out.Trades = append(out.Trades, pos.trade)
				if amb {
					out.AmbiguousTrades++
				}
				pos = nil
			}
		}
		if i+1 == len(bars) {
			break
		}
		s := makeSignal(bars, e20, e50, i)
		if s.dir != 0 {
			if pos != nil {
				out.IgnoredSignals++
			} else {
				pending = &s
			}
		}
	}
	out.OpenPosition = pos != nil
	return out, nil
}

func makeSignal(b []Bar, e20, e50 []float64, confirmation int) signal {
	p := confirmation - 2
	if p < Warmup || p < 2 {
		return signal{}
	}
	low, high := true, true
	for _, j := range []int{p - 2, p - 1, p + 1, p + 2} {
		if b[p].Low >= b[j].Low {
			low = false
		}
		if b[p].High <= b[j].High {
			high = false
		}
	}
	if low == high {
		return signal{}
	}
	if low && e20[p] > e50[p] && e20[confirmation] > e50[confirmation] && b[p].Low > e50[p] && b[p].Low < e20[p] {
		return signal{1, e50[confirmation]}
	}
	if high && e20[p] < e50[p] && e20[confirmation] < e50[confirmation] && b[p].High > e20[p] && b[p].High < e50[p] {
		return signal{-1, e50[confirmation]}
	}
	return signal{}
}
func enter(b Bar, s signal) *position {
	price := b.Open * (1 + float64(s.dir)*Slippage)
	risk := float64(s.dir) * (price - s.stop)
	target := price + float64(s.dir)*1.5*risk
	if !finite(risk) || risk <= 0 || !finite(target) || target <= 0 {
		return nil
	}
	return &position{trade: Trade{Direction: s.dir, EntryTime: b.Time, Entry: price, Stop: s.stop, Target: target, Risk: risk}, rawEntry: b.Open}
}

func applyFunding(p *position, cost float64, ambiguous bool) {
	p.funding += cost
	p.trade.FundingAmbiguous = p.trade.FundingAmbiguous || ambiguous
}

func fundingOutcome(charge float64, intrabarExit bool) (float64, bool) {
	if !intrabarExit {
		return charge, false
	}
	if charge > 0 {
		return charge, true
	}
	return 0, true
}

// openExit tests only the observed opening price.
func openExit(b Bar, t Trade) (float64, string, bool, bool) {
	if t.Direction == 1 {
		if b.Open <= t.Stop {
			return b.Open, "stop", false, true
		}
		if b.Open >= t.Target {
			return t.Target, "target", false, true
		}
	} else {
		if b.Open >= t.Stop {
			return b.Open, "stop", false, true
		}
		if b.Open <= t.Target {
			return t.Target, "target", false, true
		}
	}
	return 0, "", false, false
}

// Price at an opening gap is used for stop gaps; targets receive no gap improvement.
func exit(b Bar, t Trade) (float64, string, bool, bool) {
	// Opening price is observed before the remaining intrabar range.
	if t.Direction == 1 {
		if b.Open <= t.Stop {
			return b.Open, "stop", false, true
		}
		if b.Open >= t.Target {
			return t.Target, "target", false, true
		}
	} else {
		if b.Open >= t.Stop {
			return b.Open, "stop", false, true
		}
		if b.Open <= t.Target {
			return t.Target, "target", false, true
		}
	}
	stop, target := false, false
	if t.Direction == 1 {
		stop = b.Low <= t.Stop
		target = b.High >= t.Target
	} else {
		stop = b.High >= t.Stop
		target = b.Low <= t.Target
	}
	if stop {
		raw := t.Stop
		if (t.Direction == 1 && b.Open < t.Stop) || (t.Direction == -1 && b.Open > t.Stop) {
			raw = b.Open
		}
		return raw, "stop", target, true
	}
	if target {
		return t.Target, "target", false, true
	}
	return 0, "", false, false
}
func finish(t *Trade, rawEntry, rawExit, funding float64, when int64, reason string, amb bool) {
	t.Exit = rawExit * (1 - float64(t.Direction)*Slippage)
	t.ExitTime = when
	t.ExitReason = reason
	t.Ambiguous = amb
	t.FundingCost = funding
	t.Fees = Fee * (t.Entry + t.Exit)
	t.SlippageCost = Slippage * (rawEntry + rawExit)
	t.GrossR = float64(t.Direction) * (rawExit - rawEntry) / t.Risk
	t.NetR = (float64(t.Direction)*(t.Exit-t.Entry) - t.Fees - funding) / t.Risk
}
