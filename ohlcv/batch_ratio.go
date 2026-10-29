package ohlcv

import (
	"fmt"
	"math"
)

// BatchRatio estimates the trade-weighted mean and five-batch t interval.
// Statistical adequacy gates and terminal-position checks are the caller's duty.
func BatchRatio(days []Daily) (float64, Interval95, error) {
	if len(days) != 365 {
		return 0, Interval95{}, fmt.Errorf("need exactly 365 days")
	}
	return batchRatio(days, 73)
}

// DevelopmentBatchRatio uses the first 360 days only; never use for holdout.
func DevelopmentBatchRatio(days []Daily) (float64, Interval95, error) {
	if len(days) != 360 {
		return 0, Interval95{}, fmt.Errorf("need exactly 360 development days")
	}
	return batchRatio(days, 72)
}
func batchRatio(days []Daily, block int) (float64, Interval95, error) {
	var total float64
	var count int64
	for _, d := range days {
		if d.Count < 0 || !finite(d.Net) || (d.Count == 0 && d.Net != 0) {
			return 0, Interval95{}, fmt.Errorf("invalid daily observation")
		}
		if int64(d.Count) > math.MaxInt64-count {
			return 0, Interval95{}, fmt.Errorf("count overflow")
		}
		count += int64(d.Count)
		total += d.Net
	}
	if count <= 0 || !finite(total) {
		return 0, Interval95{}, fmt.Errorf("invalid totals")
	}
	theta := total / float64(count)
	var q float64
	for j := 0; j < 5; j++ {
		var z float64
		for _, d := range days[j*block : (j+1)*block] {
			z += d.Net - theta*float64(d.Count)
		}
		q += z * z
	}
	se := math.Sqrt(1.25*q) / float64(count)
	if !finite(theta) || !finite(se) || se <= 0 {
		return 0, Interval95{}, fmt.Errorf("invalid standard error")
	}
	lower, upper := theta-2.7764451051977987*se, theta+2.7764451051977987*se
	if !finite(lower) || !finite(upper) {
		return 0, Interval95{}, fmt.Errorf("invalid confidence interval")
	}
	return theta, Interval95{lower, upper}, nil
}

// BatchAdequacy applies sample-count and terminal-position gates. An error
// denotes invalid input; false denotes an otherwise valid but inadequate sample.
func BatchAdequacy(days []Daily, openPosition bool) (bool, error) {
	if len(days) != 365 {
		return false, fmt.Errorf("need exactly 365 days")
	}
	total, active := int64(0), 0
	emptyBlock := false
	for j := 0; j < 5; j++ {
		var block int64
		for _, d := range days[j*73 : (j+1)*73] {
			if d.Count < 0 || !finite(d.Net) || (d.Count == 0 && d.Net != 0) {
				return false, fmt.Errorf("invalid daily observation")
			}
			if int64(d.Count) > math.MaxInt64-total {
				return false, fmt.Errorf("count overflow")
			}
			total += int64(d.Count)
			block += int64(d.Count)
			if d.Count > 0 {
				active++
			}
		}
		if block == 0 {
			emptyBlock = true
		}
	}
	return total >= 200 && active >= 30 && !openPosition && !emptyBlock, nil
}
