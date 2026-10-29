package main

import (
	"fmt"
	"math"
	"sort"
)

type rng struct {
	s      uint64
	spare  float64
	cached bool
}

func (r *rng) u() uint64 {
	r.s += 0x9e3779b97f4a7c15
	z := r.s
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
func (r *rng) norm() float64 {
	if r.cached {
		r.cached = false
		return r.spare
	}
	a := float64(r.u()>>11) / 9007199254740992.0
	for a == 0 {
		a = float64(r.u()>>11) / 9007199254740992.0
	}
	b := float64(r.u()>>11) / 9007199254740992.0
	v := math.Sqrt(-2 * math.Log(a))
	r.spare = v * math.Sin(2*math.Pi*b)
	r.cached = true
	return v * math.Cos(2*math.Pi*b)
}
func (r *rng) index(n int) int {
	x := uint64(n)
	t := -x % x
	v := r.u()
	for v < t {
		v = r.u()
	}
	return int(v % x)
}
func main() {
	const n = 365
	const reps = 400
	const boots = 2000
	lens := []int{3, 7, 14, 28}
	for _, c := range []struct {
		name     string
		phi      float64
		variable bool
	}{{"iid", 0, false}, {"ar0.5", .5, false}, {"ar0.8", .8, false}, {"ar0.8-variable-count", .8, true}} {
		rngGen := rng{s: 0x1122334455667788}
		cover := make([]int, len(lens))
		widths := make([]float64, len(lens))
		for rep := 0; rep < reps; rep++ {
			sum := make([]float64, n)
			count := make([]float64, n)
			x := rngGen.norm()
			for d := 0; d < n; d++ {
				x = c.phi*x + math.Sqrt(1-c.phi*c.phi)*rngGen.norm()
				k := 1
				if c.variable {
					k = rngGen.index(5)
				}
				count[d] = float64(k)
				sum[d] = float64(k) * x
			}
			for j, L := range lens {
				rngBoot := rng{s: uint64(rep+1)*0x9e3779b97f4a7c15 + uint64(L)}
				vals := make([]float64, boots)
				prefix := make([]float64, 2*n+1)
				cnt := make([]float64, 2*n+1)
				for q := 0; q < 2*n; q++ {
					prefix[q+1] = prefix[q] + sum[q%n]
					cnt[q+1] = cnt[q] + count[q%n]
				}
				for b := 0; b < boots; b++ {
					var ns, ds float64
					for remain := n; remain > 0; {
						q := L
						if q > remain {
							q = remain
						}
						i := rngBoot.index(n)
						ns += prefix[i+q] - prefix[i]
						ds += cnt[i+q] - cnt[i]
						remain -= q
					}
					vals[b] = ns / ds
				}
				sort.Float64s(vals)
				lo, hi := vals[49], vals[1949]
				if lo <= 0 && hi >= 0 {
					cover[j]++
				}
				widths[j] += hi - lo
			}
		}
		fmt.Print(c.name)
		for j, L := range lens {
			fmt.Printf(" L%d coverage=%d/%d=%.3f meanwidth=%.4f", L, cover[j], reps, float64(cover[j])/reps, widths[j]/reps)
		}
		fmt.Println()
	}
}
