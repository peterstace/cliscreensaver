package main

import (
	"os"
	"time"
)

func main() {
	s := newScreen(termSize())
	for i := range s.data {
		r, c := s.rowCol(i)
		s.data[i].r = float64(r) / float64(s.rows)
		s.data[i].g = 0x5f / float64(0xff)
		s.data[i].b = float64(c) / float64(s.cols)
	}
	s.print(os.Stdout)

	time.Sleep(time.Hour)
}
