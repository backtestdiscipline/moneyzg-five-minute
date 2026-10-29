package study

import (
	"moneyzg-five-minute/ohlcv"
	"testing"
)

func TestSummaryGapStopUsesExactModeledStopFill(t *testing.T) {
	stop := 123.456
	regular := stop * (1 - ohlcv.Slippage)
	gap := 120 * (1 - ohlcv.Slippage)
	s := summarize(ohlcv.Result{Trades: []ohlcv.Trade{{Direction: 1, ExitReason: "stop", Stop: stop, Exit: regular}, {Direction: 1, ExitReason: "stop", Stop: stop, Exit: gap}}})
	if s.Stops != 2 || s.GapStops != 1 {
		t.Fatal(s)
	}
}
