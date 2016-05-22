package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/leocassarani/tinyrenderer/wavefront"
)

var (
	black = color.NRGBA{0, 0, 0, 255}
	white = color.NRGBA{255, 255, 255, 255}
	red   = color.NRGBA{255, 0, 0, 255}
	green = color.NRGBA{0, 255, 0, 255}
	blue  = color.NRGBA{0, 0, 255, 255}
)

const (
	width  = 800
	height = 800
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tinyrenderer [output file]")
		os.Exit(1)
	}

	out := os.Args[1]

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, black)
		}
	}

	renderModel("african_head.obj", img)

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

func renderModel(filename string, img *image.NRGBA) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	model, err := wavefront.ParseModel(file)
	if err != nil {
		log.Fatal(err)
	}

	for _, face := range model.Faces {
		for i := 0; i < 3; i++ {
			v0 := model.VertexAt(face.Indices[i])
			v1 := model.VertexAt(face.Indices[(i+1)%3])

			x0 := int((v0.X + 1) * float64(width) / 2)
			y0 := int((v0.Y + 1) * float64(height) / 2)

			x1 := int((v1.X + 1) * float64(width) / 2)
			y1 := int((v1.Y + 1) * float64(height) / 2)

			line(x0, y0, x1, y1, img, white)
		}
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
