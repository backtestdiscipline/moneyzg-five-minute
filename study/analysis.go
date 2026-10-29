// Package study owns the MoneyZG holdout science, not platform integrity.
package study

import (
	"bytes"
	"encoding/json"
	"fmt"
	"moneyzg-five-minute/ohlcv"
	"research-platform/seal"
	"research-platform/spec"
	"time"
)

const Start int64 = 1735689600000
const End int64 = 1767225600000

type Analysis struct {
	Version  string               `json:"version"`
	Start    int64                `json:"start_ms"`
	End      int64                `json:"end_ms"`
	Outcome  seal.TerminalOutcome `json:"outcome"`
	Verdict  string               `json:"verdict"`
	Mean     *float64             `json:"mean_net_r"`
	Interval *ohlcv.Interval95    `json:"interval_95"`
	Daily    []ohlcv.Daily        `json:"daily"`
	Ledger   ohlcv.Result         `json:"ledger"`
	Summary  Summary              `json:"summary"`
}

func Analyze(bars []ohlcv.Bar, funding []ohlcv.Funding) (Analysis, error) {
	a := Analysis{Version: "moneyzg-holdout/v1", Start: Start, End: End, Daily: make([]ohlcv.Daily, 365)}
	result, err := ohlcv.Evaluate(bars, funding, Start, End)
	if err != nil {
		return a, err
	}
	a.Ledger = result
	a.Summary = summarize(result)
	for _, t := range result.Trades {
		d := int((t.ExitTime - Start) / int64(24*time.Hour/time.Millisecond))
		if d < 0 || d >= 365 {
			return a, fmt.Errorf("trade exit outside holdout")
		}
		a.Daily[d].Net += t.NetR
		a.Daily[d].Count++
	}
	adequate, err := ohlcv.BatchAdequacy(a.Daily, result.OpenPosition)
	if err != nil {
		return a, err
	}
	if !adequate {
		a.Outcome = seal.Incomplete
		a.Verdict = "INCONCLUSIVE_INADEQUATE"
		return a, nil
	}
	mean, interval, err := ohlcv.BatchRatio(a.Daily)
	if err != nil {
		return a, err
	}
	a.Outcome = seal.Complete
	a.Mean = &mean
	a.Interval = &interval
	if interval.Low > 0 {
		a.Verdict = "MET"
	} else if interval.High < 0 {
		a.Verdict = "NEGATIVE_EXPECTANCY_SUPPORTED"
	} else {
		a.Verdict = "INCONCLUSIVE"
	}
	return a, nil
}

type Evaluation struct {
	Outcome seal.TerminalOutcome
	Payload []byte
}
type Verification struct {
	Status string `json:"status"`
	Notes  string `json:"notes"`
}

func Evaluate(in Inputs) (Evaluation, error) {
	bars, funding, e := load(in)
	if e != nil {
		return technical(e)
	}
	a, e := Analyze(bars, funding)
	if e != nil {
		return technical(e)
	}
	b, e := spec.CanonicalJSON(a)
	return Evaluation{a.Outcome, b}, e
}
func technical(err error) (Evaluation, error) {
	b, e := spec.CanonicalJSON(map[string]string{"status": "TECHNICAL_FAILURE", "reason": err.Error()})
	return Evaluation{seal.TechnicalFailure, b}, e
}

// DecodePayload rejects noncanonical, duplicate and unknown fields by re-encoding.
func DecodePayload(b []byte) (Analysis, error) {
	var a Analysis
	e := json.Unmarshal(b, &a)
	if e != nil {
		return a, e
	}
	canonical, e := spec.CanonicalJSON(a)
	if e != nil {
		return a, e
	}
	if !bytes.Equal(canonical, b) {
		return a, fmt.Errorf("noncanonical or unknown result field")
	}
	return a, nil
}
