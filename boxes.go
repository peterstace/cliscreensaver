package main

import (
	"fmt"
	"io"
	"syscall"
	"unsafe"
)

func newScreen(rows, cols int) screen {
	return screen{
		rows: rows,
		cols: cols,
		data: make([]float64, rows*cols),
	}
}

type screen struct {
	rows int
	cols int
	data []float64 // row-major order
}

func (s *screen) rowCol(i int) (int, int) {
	return i / s.cols, i % s.cols
}

func (s *screen) print(w io.Writer) {
	writeCursorPosition(w, 1, 1)
	w.Write([]byte("\x1b[37;40m"))
	for i := range s.data {
		level, qerr := quantize(s.data[i])
		switch level {
		case 0:
			w.Write([]byte(" "))
		case 1:
			w.Write([]byte("░"))
		case 2:
			w.Write([]byte("▒"))
		case 3:
			w.Write([]byte("▓"))
		case 4:
			w.Write([]byte("\x1b[30;47m \x1b[37;40m"))
		default:
			panic(level)
		}

		// TODO: Maybe do serpentine left/right.
		r, c := s.rowCol(i)
		qerr /= 16
		if c != s.cols-1 {
			s.data[i+1] += qerr * 7
		}
		if r != s.rows-1 {
			if c != 0 {
				s.data[i+s.cols-1] += qerr * 3
			}
			s.data[i+s.cols] += qerr * 5
			if c != s.cols-1 {
				s.data[i+s.cols+1] += qerr
			}
		}

	}
	w.Write([]byte("\x1b[m"))
}

func quantize(v float64) (int, float64) {
	switch {
	case v < 0.2:
		return 0, v - 0.1
	case v < 0.4:
		return 1, v - 0.3
	case v < 0.6:
		return 2, v - 0.5
	case v < 0.8:
		return 3, v - 0.7
	default:
		return 4, v - 0.9
	}
}

func termSize() (rows, cols int) {
	var winsize struct {
		rows, cols       uint16
		xpixels, ypixels uint16
	}
	ret, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&winsize)))

	if int(ret) == -1 {
		panic(errno)
	}
	rows = int(winsize.rows)
	cols = int(winsize.cols)
	return
}

func writeCursorPosition(w io.Writer, row, col int) {
	_, err := fmt.Fprintf(w, "\x1b[%d;%dH", row, col)
	if err != nil {
		panic(err)
	}
}
