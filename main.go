package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var (
	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tinyrenderer [output file]")
		os.Exit(1)
	}

	out := os.Args[1]

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, black)
		}
	}

	line(13, 20, 80, 40, img, white)
	line(20, 13, 40, 80, img, red)
	line(80, 40, 13, 20, img, red)

	file, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0655)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		log.Fatal(err)
	}
}

func line(x0, y0, x1, y1 int, img *image.RGBA, color color.RGBA) {
	var x, y int
	for t := 0.0; t < 1; t += 0.01 {
		x = int(float64(x0)*(1.0-t) + float64(x1)*t)
		y = int(float64(y0)*(1.0-t) + float64(y1)*t)
		img.Set(x, y, color)
	}
}
