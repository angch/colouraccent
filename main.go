package main

// https://forum.lowyat.net/index.php?act=ST&f=12&t=4908702&st=0

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"

	gocolor "github.com/gerow/go-color"
)

func main() {
	// cf https://pictr.com/image/5FrO4v
	inFile := "5FrO4v.png"
	r, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	
	img, err := jpeg.Decode(r)
	if err != nil {
		switch err.(type) {
		case jpeg.FormatError:
			// Not a jpeg, retrying as PNG
			r.Seek(0,0)
			img, err = png.Decode(r)	
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	// log.Printf("%+v\n", img)
	bound := img.Bounds()
	// log.Printf("%d %d\n", bound.Max.X, bound.Max.Y)

	counts := make(map[gocolor.HSL]float64)
	for y := 0; y < bound.Max.Y; y++ {
		for x := 0; x < bound.Max.X; x++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()
			if a < 65535 {
				continue // Not doing alpha
			}
			rgb := gocolor.RGB{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535.0}
			hsv := rgb.ToHSL()
			// log.Printf("%d %d %d %d", r, g, b, a)
			// log.Printf(" %d %d %d %d\n", hsv.H, hsv.L, hsv.S)
			if hsv.L < 0.0001 {
				continue
			}

			hsv.S = math.Round(hsv.S*256) / 256
			hsv.H = math.Round(hsv.H*256) / 256
			intensity := hsv.L
			hsv.L = 0

			// Reject boring grays
			if hsv.S < 0.3 {
				continue
			}

			counts[hsv] += intensity

			// log.Printf("%d %d %d %d\n", r, g, b, a)
		}
	}

	// Find brightest
	// log.Println(counts)
	brigtestHSV := gocolor.HSL{}
	for k, v := range counts {
		// Take grayness (aka saturation) into account too.
		if brigtestHSV.L < v*k.S {
			// Yeah, need copy, not reference a moving struct
			brigtestHSV.H = k.H
			brigtestHSV.S = k.S
			brigtestHSV.L = v * k.S
			// copy(brightestHSV, k)
		}
	}

	colors := ""
	for l := 0.0; l < 1; l += 0.1 {
		brigtestHSV.L = l
		rgb := brigtestHSV.ToRGB()
		colors += fmt.Sprintf("<div%%20style='background:%%23%02x%02x%02x;'>Accent</div>", int(rgb.R*255), int(rgb.G*255), int(rgb.B*255))
	}
	fmt.Println("Copy paste this to browser:")
	fmt.Printf("data:text/html,<html><body>%s</body></html>\n", colors)
}
