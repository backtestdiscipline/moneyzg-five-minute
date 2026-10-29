package study

import "moneyzg-five-minute/ohlcv"

// Summary is descriptive; a nil ratio denotes an undefined statistic.
type Summary struct {
	Trades, Longs, Shorts, Wins, Stops, Targets, GapStops                   int
	AmbiguousTrades, FundingAmbiguousTrades                                 int
	GrossExpectancy, NetExpectancy, WinRate, ProfitFactor, AmbiguousPercent *float64
	MaxDrawdown                                                             *float64
	Fees, Slippage, Funding                                                 float64
}

func summarize(result ohlcv.Result) Summary {
	var s Summary
	s.Trades = len(result.Trades)
	s.AmbiguousTrades = result.AmbiguousTrades
	var gross, net, up, down float64
	wealth, peak, drawdown := 100.0, 100.0, 0.0
	insolvent := false
	for _, t := range result.Trades {
		if t.Direction == 1 {
			s.Longs++
		} else {
			s.Shorts++
		}
		if t.NetR > 0 {
			s.Wins++
			up += t.NetR
		} else if t.NetR < 0 {
			down -= t.NetR
		}
		if t.ExitReason == "stop" {
			s.Stops++
			stopFill := t.Stop * (1 - float64(t.Direction)*ohlcv.Slippage)
			if (t.Direction == 1 && t.Exit < stopFill) || (t.Direction == -1 && t.Exit > stopFill) {
				s.GapStops++
			}
		} else if t.ExitReason == "target" {
			s.Targets++
		}
		if t.FundingAmbiguous {
			s.FundingAmbiguousTrades++
		}
		gross += t.GrossR
		net += t.NetR
		s.Fees += t.Fees
		s.Slippage += t.SlippageCost
		s.Funding += t.FundingCost
		wealth += t.NetR
		if wealth <= 0 {
			insolvent = true
		}
		if wealth > peak {
			peak = wealth
		}
		if peak > 0 && (peak-wealth)/peak > drawdown {
			drawdown = (peak - wealth) / peak
		}
	}
	if s.Trades > 0 {
		n := float64(s.Trades)
		g, v, w, a := gross/n, net/n, float64(s.Wins)/n, 100*float64(s.AmbiguousTrades)/n
		s.GrossExpectancy = &g
		s.NetExpectancy = &v
		s.WinRate = &w
		s.AmbiguousPercent = &a
	}
	if down > 0 {
		p := up / down
		s.ProfitFactor = &p
	}
	if !insolvent {
		s.MaxDrawdown = &drawdown
	}
	return s
}
