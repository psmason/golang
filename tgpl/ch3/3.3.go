package main

import (
	"errors"
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			var ax, ay, bx, by, cx, cy, dx, dy float64
			var err error
			if ax, ay, err = corner(i+1, j); err != nil {
				continue
			}
			if bx, by, err = corner(i, j); err != nil {
				continue
			}
			if cx, cy, err = corner(i, j+1); err != nil {
				continue
			}
			if dx, dy, err = corner(i+1, j+1); err != nil {
				continue
			}

			strokeColor := stroke(i, j)
			fmt.Printf("<polygon points='%g,%g,%g,%g,%g,%g,%g,%g' stroke='%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, strokeColor)
		}
	}
	fmt.Println("</svg>")
}

func stroke(i, j int) string {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)

	red := int(255 * 0.5 * (1 - z))
	green := int(0)
	blue := int(255 * 0.5 * (1 + z))

	return fmt.Sprintf("#%02x%02x%02x", red, green, blue)
}

func corner(i, j int) (float64, float64, error) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsInf(z, 1) {
		fmt.Fprintf(os.Stderr, "infinite value found\n")
		return -1, -1, errors.New("infinite value")
	}

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, nil
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
