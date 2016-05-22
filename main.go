package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
)

var (
	black = color.NRGBA{0, 0, 0, 255}
	white = color.NRGBA{255, 255, 255, 255}
	red   = color.NRGBA{255, 0, 0, 255}
	green = color.NRGBA{0, 255, 0, 255}
	blue  = color.NRGBA{0, 0, 255, 255}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tinyrenderer [output file]")
		os.Exit(1)
	}

	out := os.Args[1]

	img := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, black)
		}
	}

	line(13, 20, 80, 40, img, white)
	line(20, 13, 40, 80, img, red)
	line(80, 40, 13, 20, img, red)

	// Flip vertically so the origin is in the bottom-left corner.
	img = imaging.FlipV(img)

	file, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0655)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		log.Fatal(err)
	}
}

func line(x0, y0, x1, y1 int, img *image.NRGBA, color color.NRGBA) {
	steep := abs(x0-x1) < abs(y0-y1)

	// If the line is steep, swap x and y.
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}

	if x0 > x1 {
		// Make it left-to-right.
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	dx := x1 - x0
	dy := y1 - y0

	derr := abs(dy) * 2
	err := 0

	y := y0

	for x := x0; x <= x1; x++ {
		if steep {
			img.Set(y, x, color)
		} else {
			img.Set(x, y, color)
		}

		err += derr
		if err > dx {
			if y1 > y0 {
				y++
			} else {
				y--
			}
			err -= dx * 2
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
