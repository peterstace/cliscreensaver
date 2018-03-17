package main

type rgb struct {
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

func foo(c rgb) (code int, residual rgb) {
	r := bound(int(c.r * 6))
	g := bound(int(c.g * 6))
	b := bound(int(c.b * 6))

	code = 16 + b + 6*g + 36*r
	residual = rgb{
		r: c.r - float64(r)/5.0,
		g: c.g - float64(g)/5.0,
		b: c.b - float64(b)/5.0,
	}
	return
}
