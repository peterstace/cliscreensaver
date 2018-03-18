package main

import (
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	c1 := randColor()
	c2 := randColor()

	s := newScreen(termSize())
	for {
		c1 = capp(add(c1, scale(randColor(), 0.3*(rand.Float64()-0.5))))
		c2 = capp(add(c2, scale(randColor(), 0.3*(rand.Float64()-0.5))))

		for i := range s.data {
			r, c := s.rowCol(i)
			f := (float64(r) + float64(c)) / float64(s.rows+s.cols)
			s.data[i] = lerp(c1, c2, f)
		}
		s.print(os.Stdout)
		time.Sleep(time.Second / 30)
	}
}
