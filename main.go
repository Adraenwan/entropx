package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"math"
)

var Usage = func() {
	fmt.Printf("usage : entropx [options...] input output\n\n")
	fmt.Printf("options :\n")
	flag.PrintDefaults()
}

func main() {
	imgColumn := flag.Int("col", 0, "image column number. 0 means automatic")
	colorPalette := flag.String("palette", "bytecode", "color palette")
	curveGenerator := flag.String("curve", "sweep", "visualization curve sweep|zigzag|zorder")

	flag.Usage = Usage
	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	// read file
	file, err := os.OpenFile(flag.Args()[0], os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	fileSize := int(fileInfo.Size())
	if *imgColumn == 0 {
		n := 1
		for n*2 < int(math.Sqrt(float64(fileSize))) {
			n *= 2
		}
		*imgColumn = n
	}
	//imgRow := 1 + int(fileSize) / *imgColumn
	imgRow := 0
	for imgRow*(*imgColumn) <= fileSize {
		imgRow += *imgColumn
	}

	in := make(chan byte, 100)
	go func() {
		buff := make([]byte, 1)
		for {
			_, err := file.Read(buff)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			in <- buff[0]
		}
		file.Close()
		close(in)
	}()

	colors, ok := ColorConverter(*colorPalette, in)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: palette \"%s\" does not exists.", *colorPalette)
		os.Exit(1)
	}

	curve, ok := Curve(*curveGenerator, *imgColumn, imgRow)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: curve \"%s\" does not exists.", *curveGenerator)
		os.Exit(1)
	}

	// build image from color and position generators
	imgRectange := image.Rectangle{image.Point{0, 0}, image.Point{*imgColumn, imgRow}}
	img := image.NewRGBA(imgRectange)

	for point := range curve {
		c, ok := <-colors
		if ok {
			img.SetRGBA(point.X, point.Y, c)
		} else {
			break
		}
	}

	fileOut, err := os.OpenFile(flag.Args()[1], os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	err = png.Encode(fileOut, img)
	if err != nil {
		panic(err)
	}
}
