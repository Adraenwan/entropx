package main

import (
	"image"
	"math"
)

func Curve(name string, column, row int) (<-chan image.Point, bool) {
	out := make(chan image.Point, 100)

	switch name {
	case "sweep":
		go CurveSweep(out, column, row)

	case "zigzag":
		go CurveZigzag(out, column, row)

	case "zorder":
		go CurveZOrder(out, column, row)

	default:
		return nil, false
	}

	return out, true
}

func CurveSweep(out chan<- image.Point, column, row int) {
	for y := 0; y < row; y++ {
		for x := 0; x < column; x++ {
			out <- image.Point{x, y}
		}
	}

	close(out)
}

func CurveZigzag(out chan<- image.Point, column, row int) {
	x, y, adder := 0, 0, 0

	for y <= row {
		out <- image.Point{x, y}

		x += adder
		if x == 0 {
			adder = 1
			y++
		} else if x == column-1 {
			adder = -1
			y++
		}
	}

	close(out)
}

func CurveZOrder(out chan<- image.Point, column, row int) {
	level := math.Ilogb(float64(column))
	for i := 0; i <= row/column; i++ {
		zorder(level, column/2, i*column, out, false)
	}
	close(out)
}

func zorder(level, size, offset int, out chan<- image.Point, closeChan bool) {
	if level == 1 {
		out <- image.Point{0, offset}
		out <- image.Point{1, offset}
		out <- image.Point{0, offset + 1}
		out <- image.Point{1, offset + 1}
	} else {
		in := make(chan image.Point, 10)
		go zorder(level-1, size/2, offset, in, true)
		for p := range in {
			out <- p
		}

		in = make(chan image.Point, 10)
		go zorder(level-1, size/2, offset, in, true)
		for p := range in {
			out <- image.Point{p.X + size, p.Y}
		}

		in = make(chan image.Point, 10)
		go zorder(level-1, size/2, offset, in, true)
		for p := range in {
			out <- image.Point{p.X, p.Y + size}
		}

		in = make(chan image.Point, 10)
		go zorder(level-1, size/2, offset, in, true)
		for p := range in {
			out <- image.Point{p.X + size, p.Y + size}
		}
	}

	if closeChan {
		close(out)
	}
}
