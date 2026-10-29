package study

import (
	"fmt"
	"math"
	"moneyzg-five-minute/ohlcv"
	"reflect"
	"research-platform/seal"
)

// independentReplay rebuilds signals ahead of execution from completed five-bar
// windows and uses a separate state machine; it never calls ohlcv.Evaluate.
func independentReplay(b []ohlcv.Bar, f []ohlcv.Funding) (ohlcv.Result, error) {
	var out ohlcv.Result
	if len(b) != 500+105120 {
		return out, fmt.Errorf("wrong bar count")
	}
	finite := func(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }
	for i, v := range b {
		if v.Time != Start+int64(i-500)*300000 || !finite(v.Open) || !finite(v.High) || !finite(v.Low) || !finite(v.Close) || !finite(v.Volume) || v.Low <= 0 || v.Volume < 0 || v.High < v.Low || v.Open < v.Low || v.Open > v.High || v.Close < v.Low || v.Close > v.High {
			return out, fmt.Errorf("invalid bar %d", i)
		}
	}
	funds := map[int64]ohlcv.Funding{}
	const step int64 = 28800000
	for _, v := range f {
		if v.Time < Start || v.Time >= End || v.Time%step > 1000 || !finite(v.Rate) || !finite(v.Mark) || v.Mark <= 0 {
			return out, fmt.Errorf("invalid funding")
		}
		key := v.Time - v.Time%step
		if _, ok := funds[key]; ok {
			return out, fmt.Errorf("duplicate funding")
		}
		funds[key] = v
	}
	for t := Start; t < End; t += step {
		if _, ok := funds[t]; !ok {
			return out, fmt.Errorf("missing funding")
		}
	}
	// Compute EMA and every eligible signal from only the pivot and its neighbours.
	e20 := make([]float64, len(b))
	e50 := make([]float64, len(b))
	e20[0] = b[0].Close
	e50[0] = b[0].Close
	for i := 1; i < len(b); i++ {
		e20[i] = e20[i-1] + 2.0/21*(b[i].Close-e20[i-1])
		e50[i] = e50[i-1] + 2.0/51*(b[i].Close-e50[i-1])
	}
	type order struct {
		direction int
		stop      float64
	}
	signals := map[int]order{}
	for p := 500; p+3 < len(b); p++ {
		low, high := true, true
		for _, j := range []int{p - 2, p - 1, p + 1, p + 2} {
			if b[p].Low >= b[j].Low {
				low = false
			}
			if b[p].High <= b[j].High {
				high = false
			}
		}
		c := p + 2
		if low && !high && e20[p] > e50[p] && e20[c] > e50[c] && b[p].Low > e50[p] && b[p].Low < e20[p] {
			signals[c] = order{1, e50[c]}
		}
		if high && !low && e20[p] < e50[p] && e20[c] < e50[c] && b[p].High > e20[p] && b[p].High < e50[p] {
			signals[c] = order{-1, e50[c]}
		}
	}
	type active struct {
		t        ohlcv.Trade
		rawEntry float64
	}
	var pos *active
	var next *order
	closeTrade := func(v ohlcv.Bar, raw float64, reason string, amb bool) {
		t := pos.t
		t.Exit = raw * (1 - float64(t.Direction)*0.0002)
		t.ExitTime = v.Time
		t.ExitReason = reason
		t.Ambiguous = amb
		t.Fees = (t.Entry + t.Exit) * 0.0005
		t.SlippageCost = (pos.rawEntry + raw) * 0.0002
		t.GrossR = float64(t.Direction) * (raw - pos.rawEntry) / t.Risk
		t.NetR = (float64(t.Direction)*(t.Exit-t.Entry) - t.Fees - t.FundingCost) / t.Risk
		out.Trades = append(out.Trades, t)
		if amb {
			out.AmbiguousTrades++
		}
		pos = nil
	}
	openCheck := func(v ohlcv.Bar) bool {
		if pos == nil {
			return false
		}
		t := pos.t
		d := float64(t.Direction)
		if d*(v.Open-t.Stop) <= 0 {
			closeTrade(v, v.Open, "stop", false)
			return true
		}
		if d*(v.Open-t.Target) >= 0 {
			closeTrade(v, t.Target, "target", false)
			return true
		}
		return false
	}
	for i := 500; i < len(b); i++ {
		v := b[i]
		f, has := funds[v.Time]
		if pos != nil && has && f.Time == v.Time {
			pos.t.FundingCost += float64(pos.t.Direction) * f.Rate * f.Mark
		}
		openCheck(v)
		if next != nil {
			d := next.direction
			entry := v.Open * (1 + float64(d)*0.0002)
			risk := float64(d) * (entry - next.stop)
			target := entry + float64(d)*1.5*risk
			if risk <= 0 || !finite(risk) || target <= 0 || !finite(target) {
				out.InvalidEntries++
			} else {
				pos = &active{t: ohlcv.Trade{Direction: d, EntryTime: v.Time, Entry: entry, Stop: next.stop, Target: target, Risk: risk}, rawEntry: v.Open}
			}
			next = nil
		}
		openCheck(v)
		if pos != nil {
			t := pos.t
			stop, target := false, false
			if t.Direction == 1 {
				stop = v.Low <= t.Stop
				target = v.High >= t.Target
			} else {
				stop = v.High >= t.Stop
				target = v.Low <= t.Target
			}
			if has && f.Time > v.Time {
				charge := float64(t.Direction) * f.Rate * f.Mark
				if stop || target {
					pos.t.FundingAmbiguous = true
					if charge > 0 {
						pos.t.FundingCost += charge
					}
				} else {
					pos.t.FundingCost += charge
				}
			}
			if stop {
				closeTrade(v, t.Stop, "stop", target)
			} else if target {
				closeTrade(v, t.Target, "target", false)
			}
		}
		if o, ok := signals[i]; ok {
			if pos != nil {
				out.IgnoredSignals++
			} else {
				x := o
				next = &x
			}
		}
	}
	out.OpenPosition = pos != nil
	return out, nil
}
func independentAnalyze(b []ohlcv.Bar, f []ohlcv.Funding) (Analysis, error) {
	a := Analysis{Version: "moneyzg-holdout/v1", Start: Start, End: End, Daily: make([]ohlcv.Daily, 365)}
	r, e := independentReplay(b, f)
	if e != nil {
		return a, e
	}
	a.Ledger = r
	a.Summary = summarize(r)
	for _, t := range r.Trades {
		j := int((t.ExitTime - Start) / 86400000)
		if j < 0 || j >= 365 {
			return a, fmt.Errorf("exit outside holdout")
		}
		a.Daily[j].Net += t.NetR
		a.Daily[j].Count++
	}
	count, days, empty := 0, 0, false
	var sum float64
	for k := 0; k < 5; k++ {
		bc := 0
		for _, d := range a.Daily[k*73 : (k+1)*73] {
			sum += d.Net
			count += d.Count
			bc += d.Count
			if d.Count > 0 {
				days++
			}
		}
		if bc == 0 {
			empty = true
		}
	}
	if r.OpenPosition || count < 200 || days < 30 || empty {
		a.Outcome = seal.Incomplete
		a.Verdict = "INCONCLUSIVE_INADEQUATE"
		return a, nil
	}
	mean := sum / float64(count)
	var q float64
	for k := 0; k < 5; k++ {
		var z float64
		for _, d := range a.Daily[k*73 : (k+1)*73] {
			z += d.Net - mean*float64(d.Count)
		}
		q += z * z
	}
	se := math.Sqrt(1.25*q) / float64(count)
	if se <= 0 || math.IsNaN(se) || math.IsInf(se, 0) {
		return a, fmt.Errorf("invalid standard error")
	}
	ci := ohlcv.Interval95{Low: mean - 2.7764451051977987*se, High: mean + 2.7764451051977987*se}
	a.Mean = &mean
	a.Interval = &ci
	a.Outcome = seal.Complete
	if ci.Low > 0 {
		a.Verdict = "MET"
	} else if ci.High < 0 {
		a.Verdict = "NEGATIVE_EXPECTANCY_SUPPORTED"
	} else {
		a.Verdict = "INCONCLUSIVE"
	}
	return a, nil
}
func Verify(in Inputs, payload []byte) (Verification, error) {
	want, e := DecodePayload(payload)
	if e != nil {
		return Verification{"FAIL", "invalid result payload"}, nil
	}
	bars, funding, e := load(in)
	if e != nil {
		return Verification{"FAIL", e.Error()}, nil
	}
	got, e := independentAnalyze(bars, funding)
	if e != nil {
		return Verification{"FAIL", e.Error()}, nil
	}
	if !reflect.DeepEqual(want, got) {
		return Verification{"FAIL", "independent full ledger or inference mismatch"}, nil
	}
	return Verification{"PASS", "independent signal, full trade ledger, funding and five-batch replay"}, nil
}
