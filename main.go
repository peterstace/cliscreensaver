package main

import (
	"os"
	"time"
)

func main() {
	s := newScreen(termSize())
	for i := range s.data {
		r, c := s.rowCol(i)
		s.data[i] = (float64(r) + float64(c)) / (float64(s.rows) + float64(s.cols))
	}
	s.print(os.Stdout)

	time.Sleep(10 * time.Second)
}
