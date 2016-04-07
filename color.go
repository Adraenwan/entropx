package main

import (
	"image/color"
)

func ColorConverter(name string, in <-chan byte) (<-chan color.RGBA, bool) {
	out := make(chan color.RGBA, 100)

	switch name {
	case "bytecode":
		go ColorByteCode(in, out)

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
