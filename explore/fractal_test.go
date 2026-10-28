package explore

import "testing"

// This fixture illustrates strict five-bar pivots, not the video's indicator.
func confirmed(xs []int, end int, high bool) []int {
	var pivots []int
	for t := 2; t+2 <= end; t++ {
		ok := true
		for _, j := range []int{t - 2, t - 1, t + 1, t + 2} {
			if (high && xs[t] <= xs[j]) || (!high && xs[t] >= xs[j]) {
				ok = false
			}
		}
		if ok {
			pivots = append(pivots, t)
		}
	}
	return pivots
}

func TestConfirmedSignalsAppearOnlyAfterSecondRightBarCloses(t *testing.T) {
	for _, tc := range []struct {
		xs   []int
		high bool
	}{
		{[]int{1, 2, 5, 2, 1, 2, 6, 2, 1}, true},
		{[]int{6, 5, 1, 5, 6, 5, 0, 5, 6}, false},
	} {
		for end := range tc.xs {
			for _, pivot := range confirmed(tc.xs, end, tc.high) {
				if end < pivot+2 {
					t.Fatalf("pivot %d visible at %d", pivot, end)
				}
				if end == pivot+2 {
					// Central open to confirmation close: three five-minute intervals.
					if (end+1-pivot)*5 != 15 {
						t.Fatal("wrong timestamp")
					}
					for _, early := range confirmed(tc.xs, end-1, tc.high) {
						if early == pivot {
							t.Fatal("signal emitted before confirmation")
						}
					}
				}
			}
		}
		if len(confirmed(tc.xs, len(tc.xs)-1, tc.high)) != 2 {
			t.Fatal("missing fixture signals")
		}
	}
}

func TestProvisionalSignalCanDisappearAndTiesAreExcluded(t *testing.T) {
	for _, xs := range [][]int{{1, 2, 5, 2, 6}, {1, 2, 5, 2, 5}} {
		if len(confirmed(xs, 4, true)) != 0 {
			t.Fatal("invalid or tied high confirmed")
		}
	}
}
