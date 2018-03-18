package main

import (
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
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

func main() {
	s := newScreen(termSize())
	c1 := lrgb{rand.Float64(), rand.Float64(), rand.Float64()}
	c2 := lrgb{rand.Float64(), rand.Float64(), rand.Float64()}
	for i := range s.data {
		r, c := s.rowCol(i)
		f := (float64(r) + float64(c)) / float64(s.rows+s.cols)
		s.data[i] = lerp(c1, c2, f)
	}
	s.print(os.Stdout)

	time.Sleep(time.Hour)
}
