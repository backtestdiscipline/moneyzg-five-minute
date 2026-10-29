// Synthetic coverage only; no market data are read.
package main

import (
	"fmt"
	"math"
	"moneyzg-five-minute/ohlcv"
)

type rng struct {
	s     uint64
	cache bool
	spare float64
}

func (r *rng) u() uint64 {
	r.s += 0x9e3779b97f4a7c15
	z := r.s
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
func (r *rng) normal() float64 {
	if r.cache {
		r.cache = false
		return r.spare
	}
	a := float64(r.u()>>11) / 9007199254740992.0
	for a == 0 {
		a = float64(r.u()>>11) / 9007199254740992.0
	}
	b := float64(r.u()>>11) / 9007199254740992.0
	v := math.Sqrt(-2 * math.Log(a))
	r.spare = v * math.Sin(2*math.Pi*b)
	r.cache = true
	return v * math.Cos(2*math.Pi*b)
}
func (r *rng) index(n uint64) int {
	threshold := -n % n
	x := r.u()
	for x < threshold {
		x = r.u()
	}
	return int(x % n)
}
func main() {
	const reps = 3000
	for i, name := range []string{"gaussian", "student_t3", "ar05", "ar08", "ar08_variable_count", "volatility_clustering"} {
		r := rng{s: 0xb8c2a1d41760f234 + uint64(i)}
		covered := 0
		fails := 0
		var widths float64
		for k := 0; k < reps; k++ {
			d := make([]ohlcv.Daily, 365)
			x := r.normal()
			v := r.normal()
			for j := range d {
				var y float64
				count := 1
				switch name {
				case "gaussian":
					y = r.normal()
				case "student_t3":
					a, b, c := r.normal(), r.normal(), r.normal()
					y = r.normal() / math.Sqrt((a*a+b*b+c*c)/3)
				case "ar05":
					x = .5*x + math.Sqrt(.75)*r.normal()
					y = x
				case "ar08", "ar08_variable_count":
					x = .8*x + .6*r.normal()
					y = x
					if name == "ar08_variable_count" {
						count = r.index(5)
					}
				case "volatility_clustering":
					v = .8*v + .6*r.normal()
					y = r.normal() * math.Exp(.5*v)
				}
				d[j] = ohlcv.Daily{Net: float64(count) * y, Count: count}
			}
			_, ci, e := ohlcv.BatchRatio(d)
			if e != nil {
				fails++
				continue
			}
			if ci.Low <= 0 && ci.High >= 0 {
				covered++
			}
			widths += ci.High - ci.Low
		}
		p := float64(covered) / reps
		z := 1.959963984540054
		den := 1 + z*z/reps
		mid := (p + z*z/(2*reps)) / den
		half := z * math.Sqrt(p*(1-p)/reps+z*z/(4*reps*reps)) / den
		fmt.Printf("%s covered=%d/%d coverage=%.4f wilson=[%.4f,%.4f] mcse=%.4f failures=%d mean_width=%.4f\n", name, covered, reps, p, mid-half, mid+half, math.Sqrt(p*(1-p)/reps), fails, widths/float64(reps-fails))
	}
}
