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
	for _, d := range s.data {
		switch {
		case d < 0.2:
			w.Write([]byte(" "))
		case d < 0.4:
			w.Write([]byte("░"))
		case d < 0.6:
			w.Write([]byte("▒"))
		case d < 0.8:
			w.Write([]byte("▓"))
		default:
			w.Write([]byte("\x1b[30;47m \x1b[37;40m"))
		}
	}
	w.Write([]byte("\x1b[m"))
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
