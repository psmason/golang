package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

const (
	cells   = 100
	xyrange = 30.0
	angle   = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	http.HandleFunc("/surface", surface)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func surface(w http.ResponseWriter, r *http.Request) {
	color := "grey"
	width := 600
	height := 320

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Failed to parse the form\n")
		return
	}

	for k, v := range r.Form {
		if "height" == k {
			token := v[0]
			var err error
			height, err = strconv.Atoi(token)
			if err != nil {
				fmt.Fprintf(w, "Failed to parse height input\n")
				return
			}
		} else if "width" == k {
			token := v[0]
			var err error
			width, err = strconv.Atoi(token)
			if err != nil {
				fmt.Fprintf(w, "Failed to parse width input\n")
				return
			}
		} else if "color" == k {
			color = v[0]
		}
	}

	xyscale := float64(width / 2 / xyrange)
	zscale := float64(height) * 0.4

	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: %s; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", color, width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			var ax, ay, bx, by, cx, cy, dx, dy float64
			var err error
			if ax, ay, err = corner(i+1, j, width, height, xyscale, zscale); err != nil {
				continue
			}
			if bx, by, err = corner(i, j, width, height, xyscale, zscale); err != nil {
				continue
			}
			if cx, cy, err = corner(i, j+1, width, height, xyscale, zscale); err != nil {
				continue
			}
			if dx, dy, err = corner(i+1, j+1, width, height, xyscale, zscale); err != nil {
				continue
			}

			fmt.Fprintf(w, "<polygon points='%g,%g,%g,%g,%g,%g,%g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintf(w, "</svg>")
}

func corner(i int, j int, width int, height int, xyscale float64, zscale float64) (float64, float64, error) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	if math.IsInf(z, 1) {
		fmt.Fprintf(os.Stderr, "infinite value found\n")
		return -1, -1, errors.New("infinite value")
	}

	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, nil
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
