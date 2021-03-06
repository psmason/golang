package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"strconv"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/lissajous", lissajous)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler (w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Failed to parse the form")
	}

	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

var palette = []color.Color{
	color.Black,
	color.RGBA{0x00, 0xFF, 0x00, 0xFF}, // green
}

const (
	blackIndex = 0
	greenIndex = 1
	redIndex   = 2
	blueIndex  = 3
)

func lissajous(w http.ResponseWriter, r *http.Request) {
	const (
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)

	cycles := 5.0

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Failed to parse the form\n")
		return
	}

	for k, v := range r.Form {
		if "cycles" == k {
			token := v[0]
			cyclesInt, err := strconv.Atoi(token)
			if err != nil {
				fmt.Fprintf(w, "Failed to parse cycle input\n")
				return
			}
			cycles = float64(cyclesInt)
		}
	}

	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size*0.5), greenIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(w, &anim)
}
