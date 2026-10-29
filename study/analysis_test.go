package study

import (
	"moneyzg-five-minute/ohlcv"
	"reflect"
	"research-platform/seal"
	"research-platform/spec"
	"testing"
)

func synthetic() ([]ohlcv.Bar, []ohlcv.Funding) {
	n := 500 + 105120
	b := make([]ohlcv.Bar, n)
	for i := range b {
		c := 100 + float64(i)*0.01
		b[i] = ohlcv.Bar{Time: Start + int64(i-500)*300000, Open: c, High: c + 0.02, Low: c - 0.01, Close: c, Volume: 1}
		if i >= 502 && i < n-200 && i%12 == 0 {
			b[i].Low = c - .15
		}
	}
	var f []ohlcv.Funding
	for t := Start; t < End; t += 28800000 {
		f = append(f, ohlcv.Funding{Time: t, Rate: 0.0001, Mark: 100})
	}
	return b, f
}
func TestIndependentReplayMatchesSyntheticYear(t *testing.T) {
	b, f := synthetic()
	a, e := Analyze(b, f)
	if e != nil {
		t.Fatal(e)
	}
	got, e := independentAnalyze(b, f)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(a, got) {
		t.Fatalf("independent replay mismatch: primary %s/%s replay %s/%s", a.Outcome, a.Verdict, got.Outcome, got.Verdict)
	}
	if a.Outcome != seal.Complete || a.Interval == nil {
		t.Fatalf("fixture did not exercise complete inference: %s %s", a.Outcome, a.Verdict)
	}
	if len(a.Ledger.Trades) == 0 {
		t.Fatal("fixture did not exercise trades")
	}
}
func TestIncompleteAndTechnicalFailures(t *testing.T) {
	b, f := synthetic()
	for i := range b {
		b[i].Open = 100
		b[i].Close = 100
		b[i].High = 101
		b[i].Low = 99
	}
	a, e := Analyze(b, f)
	if e != nil || a.Outcome != seal.Incomplete {
		t.Fatalf("empty valid run: %+v %v", a, e)
	}
	b[501].Time++
	if _, e = Analyze(b, f); e == nil {
		t.Fatal("invalid bar accepted")
	}
	b[501].Time--
	f = f[1:]
	if _, e = Analyze(b, f); e == nil {
		t.Fatal("missing funding accepted")
	}
}
func TestCorruptResultDoesNotMatch(t *testing.T) {
	b, f := synthetic()
	a, e := Analyze(b, f)
	if e != nil {
		t.Fatal(e)
	}
	a.Ledger.Trades[0].NetR++
	got, e := independentAnalyze(b, f)
	if e != nil {
		t.Fatal(e)
	}
	if reflect.DeepEqual(a, got) {
		t.Fatal("corrupted ledger accepted")
	}
}
func TestPayloadRejectsUnknownAndDuplicateFields(t *testing.T) {
	a := Analysis{Version: "moneyzg-holdout/v1", Start: Start, End: End, Outcome: seal.Incomplete, Verdict: "INCONCLUSIVE_INADEQUATE", Daily: make([]ohlcv.Daily, 365)}
	b, e := spec.CanonicalJSON(a)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = DecodePayload(b); e != nil {
		t.Fatal(e)
	}
	bad := append([]byte(`{"unknown":1,`), b[1:]...)
	if _, e = DecodePayload(bad); e == nil {
		t.Fatal("unknown field accepted")
	}
	dup := append([]byte(`{"version":"bad",`), b[1:]...)
	if _, e = DecodePayload(dup); e == nil {
		t.Fatal("duplicate field accepted")
	}
}
func TestFundingDecoderRejectsDuplicateAndUnknown(t *testing.T) {
	valid := `[{"symbol":"BTCUSDT","fundingTime":1735689600000,"fundingRate":"0.001","markPrice":"100","rateType":"Regular"}]`
	if _, e := decodeFunding([]byte(valid)); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []string{`[{"symbol":"BTCUSDT","fundingTime":1735689600000,"fundingTime":1735689600001,"fundingRate":"0.001","markPrice":"100"}]`, `[{"symbol":"BTCUSDT","fundingTime":1735689600000,"fundingRate":"0.001","markPrice":"100","extra":1}]`, valid + `{}`} {
		if _, e := decodeFunding([]byte(bad)); e == nil {
			t.Fatal("malformed funding accepted")
		}
	}
}
