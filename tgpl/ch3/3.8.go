package main

import (
	"image"
	"image/color"
	"image/png"
	"math/big"
	"math/rand"
	"os"
	"testing"
)

const (
	xmin, ymin, xmax, ymax = -1.0, -0.5, -0.5, 0.5
	width, height          = 1024, 1024
	iterations             = 250
	precision              = 200
)

func main() {
	// widthF := big.NewFloat(float64(width))
	// heightF := big.NewFloat(float64(width))
	widthF  := new(big.Float).SetPrec(precision).SetInt64(width)
	heightF := new(big.Float).SetPrec(precision).SetInt64(height)
	iWidthF := Inv(widthF)
	iHeightF := Inv(heightF)
	xminF := big.NewFloat(xmin)
	yminF := big.NewFloat(ymin)
	xmaxF := big.NewFloat(xmax)
	ymaxF := big.NewFloat(ymax)
	yAdj := Mul(iHeightF,Sub(ymaxF, yminF))
	xAdj := Mul(iWidthF, Sub(xmaxF, xminF))
	ySampleWidth,_ := Mul(yAdj, big.NewFloat(0.25)).Float64()
	xSampleWidth,_ := Mul(xAdj, big.NewFloat(0.25)).Float64()
	
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		pyF := big.NewFloat(float64(py))
		y := Add(Mul(pyF, yAdj), yminF)
		for px := 0; px < width; px++ {
			pxF := big.NewFloat(float64(px))
			x := Add(Mul(pxF, xAdj), xminF)
			img.Set(px, py, mandelbrotSample(x, y, xSampleWidth, ySampleWidth, 2))
			// note the coordinates increase right and down
		}
	}
	png.Encode(os.Stdout, img)
}

func computeColor(n uint8) color.Color {
	// color scheme is black (0,0,0) to sky blue (#87cefa)
	// faster the escape (fewer iterations), the closer
	// the color is to sky blue.

	if int64(n) >= iterations {
		return color.Black
	}

	const alpha = 255
	const contrast = 10

	var red, blue, green uint8
	red, blue, green = 135, 206, 250

	red = uint8(float64(red) * (1 - float64(contrast)*float64(n)/iterations))
	blue = uint8(float64(blue) * (1 - float64(contrast)*float64(n)/iterations))
	green = uint8(float64(green) * (1 - float64(contrast)*float64(n)/iterations))

	return color.RGBA{red, blue, green, alpha}
}

func mandelbrotSample(a *big.Float, b *big.Float, xSampleWidth float64, ySampleWidth float64, n uint8) color.Color {
	total := 0.0
	for i := uint8(0); i<n; i++ {
		xEps := rand.Float64() * xSampleWidth
		yEps := rand.Float64() * ySampleWidth
		xSample := Add(big.NewFloat(xEps), a)
		ySample := Add(big.NewFloat(yEps), b)
		total = total + float64(mandelbrot(xSample, ySample))
	}

	average := uint8(total/float64(n))
	return computeColor(average)
}

func mandelbrot(a, b *big.Float) uint8 {
	vA := new(big.Float).SetPrec(precision).SetFloat64(0.0)
	vB := new(big.Float).SetPrec(precision).SetFloat64(0.0)

	f2 := new(big.Float).SetPrec(precision).SetFloat64(2.0)
	f4 := new(big.Float).SetPrec(precision).SetFloat64(4.0)
	
	for n := uint8(0); n < iterations; n++ {
		//v = v*v + z
		vA,vB = Add(Sub(Mul(vA, vA), Mul(vB, vB)),a), Add(Mul(f2, Mul(vA, vB)),b)

		escape := Add(Mul(vA, vA), Mul(vB, vB))
		if Greater(escape, f4) >= 0 {
			return n
		}
	}
	return iterations
}

func Mul(x, y *big.Float) *big.Float {
	return new(big.Float).SetPrec(precision).Mul(x, y)
}
func Sub(x, y *big.Float) *big.Float {
	return new(big.Float).SetPrec(precision).Sub(x, y)
}
func Add(x, y *big.Float) *big.Float {
	return new(big.Float).SetPrec(precision).Add(x, y)
}

func Inv(x *big.Float) *big.Float {
	return new(big.Float).SetPrec(precision).Quo(big.NewFloat(1), x)
}

func Greater(x, y *big.Float) int {
	return Sub(x, y).Sign()
}

func TestInv(t *testing.T) {
	expected := big.NewFloat(0.1);
	computed := Inv(big.NewFloat(10.0));
	if expected != computed {
		t.Errorf("mismatch: %s  %s", expected.String(), computed.String())
	}
}


