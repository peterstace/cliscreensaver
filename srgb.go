package main

import (
	"math"
	"math/rand"
)

// lrgb represents a linear RGB color. Each element should be between 0 and 1
// (inclusive).
type lrgb struct {
	r, g, b float64
}

func bound(i int) int {
	if i < 0 {
		return 0
	}
	if i > 5 {
		return 5
	}
	return i
}

func quantize(c lrgb) (code int, residual lrgb) {
	// TODO: This should really take into account where the 6 barrier
	// points are.
	r := bound(int(c.r * 6))
	g := bound(int(c.g * 6))
	b := bound(int(c.b * 6))

	code = 16 + b + 6*g + 36*r
	residual = lrgb{
		r: c.r - float64(r)/5.0,
		g: c.g - float64(g)/5.0,
		b: c.b - float64(b)/5.0,
	}
	return
}

func randColor() lrgb {
	return lrgb{
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
	}
}

func lerp(c1, c2 lrgb, f float64) lrgb {
	return lrgb{
		c1.r*(1-f) + c2.r*f,
		c1.g*(1-f) + c2.g*f,
		c1.b*(1-f) + c2.b*f,
	}
}

func scale(c lrgb, f float64) lrgb {
	return lrgb{
		c.r * f,
		c.g * f,
		c.b * f,
	}
}

func capp(c lrgb) lrgb {
	return lrgb{
		math.Min(1, math.Max(0, c.r)),
		math.Min(1, math.Max(0, c.g)),
		math.Min(1, math.Max(0, c.b)),
	}
}

func add(a, b lrgb) lrgb {
	return lrgb{
		a.r + b.r,
		a.g + b.g,
		a.b + b.b,
	}
}
