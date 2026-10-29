// Command development runs exploratory diagnostics only on verified development bytes.
package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"moneyzg-five-minute/ohlcv"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type artifact struct {
	Name     string    `json:"name"`
	Bytes    int       `json:"bytes"`
	SHA      string    `json:"sha256"`
	Checksum *artifact `json:"checksum"`
}
type rate struct {
	Symbol string `json:"symbol"`
	Time   int64  `json:"fundingTime"`
	Rate   string `json:"fundingRate"`
	Mark   string `json:"markPrice"`
}

func verified(dir string, a artifact) ([]byte, error) {
	if filepath.Base(a.Name) != a.Name || len(a.SHA) != 64 {
		return nil, fmt.Errorf("invalid artifact name or hash")
	}
	b, e := os.ReadFile(filepath.Join(dir, a.Name))
	if e != nil {
		return nil, e
	}
	s := sha256.Sum256(b)
	if len(b) != a.Bytes || hex.EncodeToString(s[:]) != a.SHA {
		return nil, fmt.Errorf("identity mismatch %s", a.Name)
	}
	return b, nil
}

// archiveWindow checks the declared CSV file's UTC calendar interval.
func archiveWindow(name string) (int64, int64, error) {
	if !strings.HasPrefix(name, "BTCUSDT-5m-") || !strings.HasSuffix(name, ".zip") {
		return 0, 0, fmt.Errorf("bad archive identity")
	}
	date := strings.TrimSuffix(strings.TrimPrefix(name, "BTCUSDT-5m-"), ".zip")
	layout := "2006-01"
	if len(date) == 10 {
		layout = "2006-01-02"
	}
	t, e := time.Parse(layout, date)
	if e != nil {
		return 0, 0, e
	}
	end := t.AddDate(0, 1, 0)
	if len(date) == 10 {
		end = t.AddDate(0, 0, 1)
	}
	return t.UnixMilli(), end.UnixMilli(), nil
}

func number(s string) (float64, error) { return strconv.ParseFloat(s, 64) }
func main() {
	dir := flag.String("data", "data/development", "verified development input directory")
	ledger := flag.String("ledger", "", "optional exploratory JSON ledger path")
	flag.Parse()
	if e := run(*dir, *ledger); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run(dir, ledger string) error {
	manifest, e := os.ReadFile("provenance/development.json")
	if e != nil {
		return e
	}
	var as []artifact
	if e = json.Unmarshal(manifest, &as); e != nil {
		return e
	}
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	end := time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC).UnixMilli()
	warm := start - 500*ohlcv.Interval
	expected := map[string]bool{"BTCUSDT-5m-2023-12.zip": true, "funding-page-1.json": true, "funding-page-2.json": true, "vision-terms.md": true}
	for m := 1; m <= 11; m++ {
		expected[fmt.Sprintf("BTCUSDT-5m-2024-%02d.zip", m)] = true
	}
	for d := 1; d <= 29; d++ {
		expected[fmt.Sprintf("BTCUSDT-5m-2024-12-%02d.zip", d)] = true
	}
	var bars []ohlcv.Bar
	var funds []ohlcv.Funding
	seen := map[string]bool{}
	for _, a := range as {
		if !expected[a.Name] {
			return fmt.Errorf("unexpected source %s", a.Name)
		}
		if seen[a.Name] {
			return fmt.Errorf("duplicate manifest entry")
		}
		seen[a.Name] = true
		b, e := verified(dir, a)
		if e != nil {
			return e
		}
		if strings.HasSuffix(a.Name, ".zip") {
			if a.Checksum == nil {
				return fmt.Errorf("missing publisher checksum")
			}
			side, e := verified(dir, *a.Checksum)
			if e != nil {
				return e
			}
			fields := strings.Fields(string(side))
			if len(fields) != 2 || fields[0] != a.SHA || fields[1] != a.Name {
				return fmt.Errorf("publisher checksum mismatch")
			}
			windowStart, windowEnd, e := archiveWindow(a.Name)
			if e != nil {
				return e
			}
			z, e := zip.NewReader(strings.NewReader(string(b)), int64(len(b)))
			if e != nil {
				return e
			}
			if len(z.File) != 1 || z.File[0].Name != strings.TrimSuffix(a.Name, ".zip")+".csv" {
				return fmt.Errorf("unexpected archive shape")
			}
			f, e := z.File[0].Open()
			if e != nil {
				return e
			}
			r := csv.NewReader(f)
			for {
				row, e := r.Read()
				if e == io.EOF {
					break
				}
				if e != nil {
					f.Close()
					return e
				}
				if len(row) != 12 {
					f.Close()
					return fmt.Errorf("bad candle columns")
				}
				if row[0] == "open_time" {
					continue
				}
				ts, e := strconv.ParseInt(row[0], 10, 64)
				if e != nil {
					f.Close()
					return e
				}
				if ts < windowStart || ts >= windowEnd {
					f.Close()
					return fmt.Errorf("bar outside archive calendar interval")
				}
				closeTime, e := strconv.ParseInt(row[6], 10, 64)
				if e != nil || closeTime != ts+ohlcv.Interval-1 {
					f.Close()
					return fmt.Errorf("invalid close time")
				}
				var v [5]float64
				for i := 0; i < 5; i++ {
					v[i], e = number(row[i+1])
					if e != nil {
						f.Close()
						return e
					}
				}
				if ts >= warm && ts < end {
					bars = append(bars, ohlcv.Bar{Time: ts, Open: v[0], High: v[1], Low: v[2], Close: v[3], Volume: v[4]})
				}
			}
			if e = f.Close(); e != nil {
				return e
			}
		} else if strings.HasPrefix(a.Name, "funding-page-") {
			var rs []rate
			if e = json.Unmarshal(b, &rs); e != nil {
				return e
			}
			for _, r := range rs {
				if r.Symbol != "BTCUSDT" {
					return fmt.Errorf("wrong symbol")
				}
				rate, e := number(r.Rate)
				if e != nil {
					return e
				}
				mark, e := number(r.Mark)
				if e != nil {
					return e
				}
				funds = append(funds, ohlcv.Funding{Time: r.Time, Rate: rate, Mark: mark})
			}
		}
	}
	if len(seen) != len(expected) {
		return fmt.Errorf("missing expected source")
	}
	sort.Slice(bars, func(i, j int) bool { return bars[i].Time < bars[j].Time })
	result, e := ohlcv.Evaluate(bars, funds, start, end)
	if e != nil {
		return e
	}
	if ledger != "" {
		b, e := json.Marshal(result)
		if e != nil {
			return e
		}
		if e = os.WriteFile(ledger, append(b, '\n'), 0600); e != nil {
			return e
		}
	}
	var long, short, stops, targets, gaps, large, verylarge int
	var feeR, slipR, fundingR float64
	distances := make([]float64, 0, len(result.Trades))
	for _, t := range result.Trades {
		if t.Direction == 1 {
			long++
		} else {
			short++
		}
		if t.ExitReason == "stop" {
			stops++
			raw := t.Exit / (1 - float64(t.Direction)*ohlcv.Slippage)
			if (t.Direction == 1 && raw < t.Stop) || (t.Direction == -1 && raw > t.Stop) {
				gaps++
			}
		} else if t.ExitReason == "target" {
			targets++
		}
		f, sl, fu := t.Fees/t.Risk, t.SlippageCost/t.Risk, t.FundingCost/t.Risk
		feeR += f
		slipR += sl
		fundingR += fu
		costs := f + sl + fu
		if costs > 0.5 {
			large++
		}
		if costs > 1 {
			verylarge++
		}
		distances = append(distances, 10000*t.Risk/t.Entry)
	}
	if len(distances) > 0 {
		sort.Float64s(distances)
		q := func(p float64) float64 { return distances[int(math.Ceil(p*float64(len(distances))))-1] }
		fmt.Printf("decomposition long=%d short=%d stops=%d targets=%d stop_gaps=%d cost_over_halfR=%d cost_over_1R=%d mean_fee_R=%.8f mean_slippage_R=%.8f mean_funding_R=%.8f initial_risk_bps_p10=%.8f median=%.8f p90=%.8f\n", long, short, stops, targets, gaps, large, verylarge, feeR/float64(len(distances)), slipR/float64(len(distances)), fundingR/float64(len(distances)), q(.1), q(.5), q(.9))
	}
	var net, gross, fees, slip, funding, positive, negative float64
	wealth, peak, maxDrawdown := 100.0, 100.0, 0.0
	insolvent := false
	wins, fa := 0, 0
	days := make([]ohlcv.Daily, int((end-start)/86400000))
	for _, t := range result.Trades {
		net += t.NetR
		gross += t.GrossR
		fees += t.Fees
		slip += t.SlippageCost
		funding += t.FundingCost
		if t.NetR > 0 {
			wins++
			positive += t.NetR
		} else if t.NetR < 0 {
			negative -= t.NetR
		}
		wealth += t.NetR
		if wealth <= 0 {
			insolvent = true
		}
		if wealth > peak {
			peak = wealth
		}
		if peak > 0 && (peak-wealth)/peak > maxDrawdown {
			maxDrawdown = (peak - wealth) / peak
		}
		if t.FundingAmbiguous {
			fa++
		}
		d := int((t.ExitTime - start) / 86400000)
		days[d].Net += t.NetR
		days[d].Count++
	}
	if len(days) >= 360 {
		mean360, ci360, e360 := ohlcv.DevelopmentBatchRatio(days[:360])
		fmt.Printf("exploratory_first_360_days_batch72_mean=%.8f ci=%+v error=%v\n", mean360, ci360, e360)
		legacy360, legacyErr := ohlcv.Bootstrap(days[:360])
		fmt.Printf("exploratory_first_360_days_legacy_block7=%+v error=%v\n", legacy360, legacyErr)
	}
	ci, err := ohlcv.Bootstrap(days)
	for _, length := range []int{3, 14, 28} {
		sensitivity, e := ohlcv.BootstrapBlocks(days, length)
		fmt.Printf("descriptive_block_%d=%+v error=%v ", length, sensitivity, e)
	}
	fmt.Println()
	if insolvent {
		fmt.Print("max_drawdown=undefined_nonpositive_equity ")
	} else {
		fmt.Printf("max_drawdown=%.8f ", maxDrawdown)
	}
	if negative > 0 {
		fmt.Printf("profit_factor=%.8f ", positive/negative)
	} else {
		fmt.Print("profit_factor=undefined ")
	}
	mean := 0.0
	if len(result.Trades) > 0 {
		n := float64(len(result.Trades))
		mean = net / n
		fmt.Printf("win_rate=%.8f gross_expectancy_R=%.8f ambiguous_pct=%.8f ", float64(wins)/n, gross/n, 100*float64(result.AmbiguousTrades)/n)
	} else {
		fmt.Print("win_rate=undefined gross_expectancy_R=undefined ambiguous_pct=undefined ")
	}
	fmt.Printf("EXPLORATORY ONLY bars=%d funding=%d trades=%d wins=%d gross_R=%.8f net_R=%.8f mean_net_R=%.8f fees=%.8f slippage=%.8f funding_cost=%.8f stop_target_ambiguous=%d funding_ambiguous=%d open=%v bootstrap=%+v bootstrap_error=%v\n", len(bars), len(funds), len(result.Trades), wins, gross, net, mean, fees, slip, funding, result.AmbiguousTrades, fa, result.OpenPosition, ci, err)
	return nil
}
