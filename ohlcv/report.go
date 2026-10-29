package ohlcv

import (
	"fmt"
	"math"
	"sort"
)

// Daily contains net-R sums and counts grouped by UTC exit day, including zero days.
type Daily struct {
	Net   float64
	Count int
}
type Interval95 struct{ Low, High float64 }

// Bootstrap uses seven-day circular blocks, SplitMix64 and nearest ranks.
func Bootstrap(days []Daily) (Interval95, error) {
	n := len(days)
	if n < 7 {
		return Interval95{}, fmt.Errorf("insufficient days")
	}
	for _, d := range days {
		if d.Count < 0 || !finite(d.Net) || (d.Count == 0 && d.Net != 0) {
			return Interval95{}, fmt.Errorf("invalid daily input")
		}
	}
	state := uint64(0x6d6f6e65797a6701)
	next := func() uint64 {
		state += 0x9e3779b97f4a7c15
		z := state
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		return z ^ (z >> 31)
	}
	bound := uint64(n)
	threshold := -bound % bound
	values := make([]float64, 10000)
	for i := range values {
		var sum float64
		count := 0
		used := 0
		for used < n {
			x := next()
			for x < threshold {
				x = next()
			}
			start := int(x % bound)
			for j := 0; j < 7 && used < n; j++ {
				d := days[(start+j)%n]
				sum += d.Net
				count += d.Count
				used++
			}
		}
		if count == 0 {
			return Interval95{}, fmt.Errorf("resample has no trades")
		}
		values[i] = sum / float64(count)
	}
	sort.Float64s(values)
	if math.IsNaN(values[249]) {
		return Interval95{}, fmt.Errorf("invalid interval")
	}
	return Interval95{values[249], values[9749]}, nil
}
