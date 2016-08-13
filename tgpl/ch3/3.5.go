package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
			// note the coordinates increase right and down
		}
	}
	png.Encode(os.Stdout, img)
}

func computeColor(n uint8) color.Color {
	// color scheme is black (0,0,0) to sky blue (#87cefa)
	// faster the escape (fewer iterations), the closer
	// the color is to sky blue.

	const alpha = 255
	const contrast = 10

	var red, blue, green uint8
	red, blue, green = 135, 206, 250

	red = uint8(float64(red) * (1 - float64(contrast)*float64(n)/200))
	blue = uint8(float64(blue) * (1 - float64(contrast)*float64(n)/200))
	green = uint8(float64(green) * (1 - float64(contrast)*float64(n)/200))

	return color.RGBA{red, blue, green, alpha}
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			//return color.Gray{255 - contrast*n}
			return computeColor(n)
		}
	}
	return color.Black
}
