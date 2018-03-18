package main

import (
	"bytes"
	"fmt"
	"io"
	"syscall"
	"unsafe"
)

func newScreen(rows, cols int) screen {
	return screen{
		rows: rows,
		cols: cols,
		data: make([]lrgb, rows*cols),
	}
}

type screen struct {
	rows int
	cols int
	data []lrgb // row-major order
}

func (s *screen) rowCol(i int) (int, int) {
	return i / s.cols, i % s.cols
}

func (s *screen) print(w io.Writer) {
	var buf bytes.Buffer
	writeCursorPosition(&buf, 1, 1)
	for i := range s.data {
		code, qerr := quantize(s.data[i])
		fmt.Fprintf(&buf, "\x1b[48;5;%dm ", code)
		r, c := s.rowCol(i)
		qerr.r /= 16.0
		qerr.g /= 16.0
		qerr.b /= 16.0

		//qerr.r = 0
		//qerr.g = 0
		//qerr.b = 0

		if c != s.cols-1 {
			s.data[i+1].r += qerr.r * 7
			s.data[i+1].g += qerr.g * 7
			s.data[i+1].b += qerr.b * 7
		}

		if r != s.rows-1 {
			if c != 0 {
				s.data[i+s.cols-1].r += qerr.r * 3
				s.data[i+s.cols-1].g += qerr.g * 3
				s.data[i+s.cols-1].b += qerr.b * 3
			}

			s.data[i+s.cols].r += qerr.r * 5
			s.data[i+s.cols].g += qerr.g * 5
			s.data[i+s.cols].b += qerr.b * 5

			if c != s.cols-1 {
				s.data[i+s.cols+1].r += qerr.r
				s.data[i+s.cols+1].g += qerr.g
				s.data[i+s.cols+1].b += qerr.b
			}
		}

	}
	buf.Write([]byte("\x1b[m")) // SGR0
	io.Copy(w, &buf)
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
