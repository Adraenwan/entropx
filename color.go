package main

import (
	"image/color"
)

func ColorConverter(name string, in <-chan byte) (<-chan color.RGBA, bool) {
	out := make(chan color.RGBA, 100)

	switch name {
	case "bytecode":
		go ColorByteCode(in, out)

	case "rainbow":
		go ColorRainbow(in, out)

	default:
		return nil, false
	}

	realOut := make(chan color.RGBA, 100)
	go func() {
		for c := range out {
			realOut <- c
		}

		for {
			realOut <- color.RGBA{0, 0, 0, 0} // VOID
		}
	}()

	return realOut, true
}

func ColorByteCode(in <-chan byte, out chan<- color.RGBA) {
	for b := range in {
		switch {
		case b == 0:
			out <- color.RGBA{0, 0, 0, 255} // BLACK

		case b == 1:
			out <- color.RGBA{255, 255, 255, 255} // WHITE

		case 31 < b && b < 126:
			out <- color.RGBA{55, 126, 184, 255} // BLUE

		default:
			out <- color.RGBA{228, 26, 28, 255} // RED
		}
	}

	close(out)
}

func ColorRainbow(in <-chan byte, out chan<- color.RGBA) {
	m1 := float64(0.5)
	m2 := float64(0.5)
	m3 := float64(0.5)

	const a1 = 0.25
	const a2 = 0.005
	const a3 = 0.001

	for b := range in {
		fb := float64(b)

		m1 = a1*fb + (1 - a1)*m1
		m2 = a2*fb + (1 - a2)*m2
		m3 = a3*fb + (1 - a3)*m3

		out <- color.RGBA{byte(m1*255), byte(m2*255), byte(m3*255), 255}
	}

	close(out)
}
