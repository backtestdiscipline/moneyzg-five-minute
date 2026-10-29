package study

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"moneyzg-five-minute/ohlcv"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Source reads use the platform's bound reference and verify at EOF and close.
func readSource(in Inputs, name string) ([]byte, error) {
	for _, ref := range in.Sources {
		if ref.Path == "data/holdout/"+name {
			r, e := in.OpenSource(ref)
			if e != nil {
				return nil, e
			}
			b, e := io.ReadAll(io.LimitReader(r, 16<<20))
			c := r.Close()
			if e != nil {
				return nil, e
			}
			if c != nil {
				return nil, c
			}
			if len(b) >= 16<<20 {
				return nil, fmt.Errorf("source exceeds limit: %s", name)
			}
			return b, nil
		}
	}
	return nil, fmt.Errorf("missing holdout source %s", name)
}
func holdoutNames() []string {
	names := []string{"BTCUSDT-5m-2024-12-30.zip", "BTCUSDT-5m-2024-12-31.zip"}
	for m := 1; m <= 12; m++ {
		names = append(names, fmt.Sprintf("BTCUSDT-5m-2025-%02d.zip", m))
	}
	return names
}
func load(in Inputs) ([]ohlcv.Bar, []ohlcv.Funding, error) {
	var bars []ohlcv.Bar
	warm := Start - 500*ohlcv.Interval
	for _, name := range holdoutNames() {
		raw, e := readSource(in, name)
		if e != nil {
			return nil, nil, e
		}
		cs, e := readSource(in, name+".CHECKSUM")
		if e != nil {
			return nil, nil, e
		}
		h := sha256.Sum256(raw)
		parts := strings.Fields(string(cs))
		if len(parts) != 2 || parts[0] != fmt.Sprintf("%x", h) || parts[1] != name {
			return nil, nil, fmt.Errorf("publisher checksum mismatch: %s", name)
		}
		z, e := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
		if e != nil {
			return nil, nil, e
		}
		if len(z.File) != 1 || z.File[0].Name != strings.TrimSuffix(name, ".zip")+".csv" || z.File[0].UncompressedSize64 > 8<<20 {
			return nil, nil, fmt.Errorf("archive contents invalid: %s", name)
		}
		date := strings.TrimSuffix(strings.TrimPrefix(name, "BTCUSDT-5m-"), ".zip")
		layout := "2006-01"
		if len(date) == 10 {
			layout = "2006-01-02"
		}
		begin, e := time.Parse(layout, date)
		if e != nil {
			return nil, nil, e
		}
		end := begin.AddDate(0, 1, 0)
		if len(date) == 10 {
			end = begin.AddDate(0, 0, 1)
		}
		f, e := z.File[0].Open()
		if e != nil {
			return nil, nil, e
		}
		r := csv.NewReader(f)
		rowNumber := 0
		for {
			row, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				f.Close()
				return nil, nil, err
			}
			if len(row) != 12 {
				f.Close()
				return nil, nil, fmt.Errorf("invalid kline schema")
			}
			rowNumber++
			if row[0] == "open_time" {
				expectedHeader := []string{"open_time", "open", "high", "low", "close", "volume", "close_time", "quote_volume", "count", "taker_buy_volume", "taker_buy_quote_volume", "ignore"}
				if !slices.Equal(row, expectedHeader) {
					f.Close()
					return nil, nil, fmt.Errorf("invalid CSV header")
				}
				if rowNumber != 1 {
					f.Close()
					return nil, nil, fmt.Errorf("unexpected repeated CSV header")
				}
				continue
			}
			t, err := strconv.ParseInt(row[0], 10, 64)
			if err != nil || t < begin.UnixMilli() || t >= end.UnixMilli() {
				f.Close()
				return nil, nil, fmt.Errorf("bar outside archive interval")
			}
			closeTime, err := strconv.ParseInt(row[6], 10, 64)
			if err != nil || closeTime != t+ohlcv.Interval-1 {
				f.Close()
				return nil, nil, fmt.Errorf("invalid close time")
			}
			v := [5]float64{}
			for i := range v {
				v[i], err = strconv.ParseFloat(row[i+1], 64)
				if err != nil {
					f.Close()
					return nil, nil, err
				}
			}
			if t >= warm && t < End {
				bars = append(bars, ohlcv.Bar{Time: t, Open: v[0], High: v[1], Low: v[2], Close: v[3], Volume: v[4]})
			}
		}
		if e = f.Close(); e != nil {
			return nil, nil, e
		}
	}
	sort.Slice(bars, func(i, j int) bool { return bars[i].Time < bars[j].Time })
	var funds []ohlcv.Funding
	for i := 1; i <= 2; i++ {
		raw, e := readSource(in, fmt.Sprintf("funding-%d.json", i))
		if e != nil {
			return nil, nil, e
		}
		records, e := decodeFunding(raw)
		if e != nil {
			return nil, nil, e
		}
		if len(records) == 0 || len(records) >= 1000 {
			return nil, nil, fmt.Errorf("funding page incomplete or truncated")
		}
		funds = append(funds, records...)
	}

	for _, ref := range in.Sources {
		if strings.HasPrefix(ref.Path, "data/holdout/") && path.Base(ref.Path) == "." {
			return nil, nil, fmt.Errorf("invalid source")
		}
	}
	return bars, funds, nil
}

// decodeFunding enforces unique known keys and rejects trailing JSON.
func decodeFunding(raw []byte) ([]ohlcv.Funding, error) {
	d := json.NewDecoder(bytes.NewReader(raw))
	tok, e := d.Token()
	if e != nil || tok != json.Delim('[') {
		return nil, fmt.Errorf("invalid funding array")
	}
	var out []ohlcv.Funding
	for d.More() {
		tok, e = d.Token()
		if e != nil || tok != json.Delim('{') {
			return nil, fmt.Errorf("invalid funding object")
		}
		seen := map[string]bool{}
		var symbol, rateText, markText, rateType string
		var ts int64
		for d.More() {
			tok, e = d.Token()
			if e != nil {
				return nil, e
			}
			key, ok := tok.(string)
			if !ok || seen[key] {
				return nil, fmt.Errorf("duplicate or invalid funding key")
			}
			seen[key] = true
			switch key {
			case "symbol":
				e = d.Decode(&symbol)
			case "fundingTime":
				e = d.Decode(&ts)
			case "fundingRate":
				e = d.Decode(&rateText)
			case "markPrice":
				e = d.Decode(&markText)
			case "rateType":
				e = d.Decode(&rateType)
			default:
				return nil, fmt.Errorf("unknown funding key %s", key)
			}
			if e != nil {
				return nil, e
			}
		}
		if _, e = d.Token(); e != nil {
			return nil, e
		}
		if !seen["symbol"] || !seen["fundingTime"] || !seen["fundingRate"] || !seen["markPrice"] || symbol != "BTCUSDT" || (rateType != "" && rateType != "Regular") {
			return nil, fmt.Errorf("invalid funding record")
		}
		rate, e := strconv.ParseFloat(rateText, 64)
		if e != nil {
			return nil, e
		}
		mark, e := strconv.ParseFloat(markText, 64)
		if e != nil {
			return nil, e
		}
		out = append(out, ohlcv.Funding{Time: ts, Rate: rate, Mark: mark})
	}
	tok, e = d.Token()
	if e != nil || tok != json.Delim(']') {
		return nil, fmt.Errorf("invalid funding end")
	}
	if _, e = d.Token(); e != io.EOF {
		return nil, fmt.Errorf("trailing funding content")
	}
	return out, nil
}
