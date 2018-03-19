package main

import (
	"math"
	"math/rand"
)

const gamma = 2.2

// lrgb represents a linear RGB color. Each element should be between 0 and 1
// (inclusive).
type lrgb struct {
	r, g, b float64
}

const (
	// levels (gamma encoded)
	lvl0 = float64(0x00) / float64(0xff)
	lvl1 = float64(0x5f) / float64(0xff)
	lvl2 = float64(0x87) / float64(0xff)
	lvl3 = float64(0xaf) / float64(0xff)
	lvl4 = float64(0xd7) / float64(0xff)
	lvl5 = float64(0xff) / float64(0xff)
)

func (c lrgb) quantize() (code int, residual lrgb) {
	c = c.bound()

	// convert to gamma space
	gar := math.Pow(c.r, 1/gamma)
	gag := math.Pow(c.g, 1/gamma)
	gab := math.Pow(c.b, 1/gamma)

	// convert to output levels
	lvlr := toLevel(gar)
	lvlg := toLevel(gag)
	lvlb := toLevel(gab)

	code = 16 + lvlb + 6*lvlg + 36*lvlr
	residual = lrgb{
		r: c.r - math.Pow(fromLevel(lvlr), gamma),
		g: c.g - math.Pow(fromLevel(lvlg), gamma),
		b: c.b - math.Pow(fromLevel(lvlb), gamma),
	}
	return
}

// gav is a gamma encoded value between 0 and 1
func toLevel(gav float64) int {
	switch {
	case gav < 0.5*lvl1:
		return 0
	case gav < 0.5*lvl1+0.5*lvl2:
		return 1
	case gav < 0.5*lvl2+0.5*lvl3:
		return 2
	case gav < 0.5*lvl3+0.5*lvl4:
		return 3
	case gav < 0.5*lvl4+0.5*lvl5:
		return 4
	default:
		return 5
	}
}

func fromLevel(lvl int) float64 {
	switch lvl {
	case 0:
		return lvl0
	case 1:
		return lvl1
	case 2:
		return lvl2
	case 3:
		return lvl3
	case 4:
		return lvl4
	case 5:
		return lvl5
	default:
		panic(lvl)
	}
}

func randColor() lrgb {
	return lrgb{
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
	}
}

func (c1 lrgb) lerp(c2 lrgb, f float64) lrgb {
	return lrgb{
		c1.r*f + c2.r*(1-f),
		c1.g*f + c2.g*(1-f),
		c1.b*f + c2.b*(1-f),
	}
}

func (c lrgb) scale(f float64) lrgb {
	return lrgb{
		c.r * f,
		c.g * f,
		c.b * f,
	}
}

func (c lrgb) bound() lrgb {
	return lrgb{
		math.Min(1, math.Max(0, c.r)),
		math.Min(1, math.Max(0, c.g)),
		math.Min(1, math.Max(0, c.b)),
	}
}

func (c lrgb) add(o lrgb) lrgb {
	return lrgb{
		c.r + o.r,
		c.g + o.g,
		c.b + o.b,
	}
}
