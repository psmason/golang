package main

import (
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = 0, -0.5, 0.5, 0.0
		width, height          = 1024, 1024
	)

	widthF := big.NewFloat(float64(width))
	heightF := big.NewFloat(float64(width))
	iWidthF := Inv(widthF)
	iHeightF := Inv(heightF)
	xminF := big.NewFloat(xmin)
	yminF := big.NewFloat(ymin)
	xmaxF := big.NewFloat(xmax)
	ymaxF := big.NewFloat(ymax)
	yDiff := Sub(ymaxF, yminF)
	xDiff := Sub(xmaxF, xminF)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		pyF := big.NewFloat(float64(py))
		y := Add(Mul(Mul(pyF, iHeightF), yDiff), yminF)
		for px := 0; px < width; px++ {
			pxF := big.NewFloat(float64(px))
			x := Add(Mul(Mul(pxF, iWidthF), xDiff), xminF)
			img.Set(px, py, mandelbrot(x, y))
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

func mandelbrot(a, b *big.Float) color.Color {
	const iterations = 200
	const contrast = 15

	vA := big.NewFloat(0)
	vB := big.NewFloat(0)
	f2 := big.NewFloat(2)
	f4 := big.NewFloat(4)
	for n := uint8(0); n < iterations; n++ {
		//v = v*v + z
		vA = Add(Sub(Mul(vA, vA), Mul(vB, vB)), a)
		vB = Add(Mul(Mul(f2, vA), vB), b)

		escape := Add(Mul(vA, vA), Mul(vB, vB))
		if Greater(escape, f4) >= 0 {
			//return color.Gray{255 - contrast*n}
			return computeColor(n)
		}
	}
	return color.Black
}

func Mul(x, y *big.Float) *big.Float {
	return big.NewFloat(0).Mul(x, y)
}
func Sub(x, y *big.Float) *big.Float {
	return big.NewFloat(0).Sub(x, y)
}
func Add(x, y *big.Float) *big.Float {
	return big.NewFloat(0).Add(x, y)
}

func Inv(x *big.Float) *big.Float {
	return big.NewFloat(0).Quo(big.NewFloat(1), x)
}

func Greater(x, y *big.Float) int {
	return Sub(x, y).Sign()
}
