package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"log"
	"math/cmplx"
	"os"
)

func mandelbrotGreyscale(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}

func mandelbrotColored(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			tmp := 255 % (n + 1)
			return palette.Plan9[tmp]
		}
	}
	return color.Black
}

func generateMandelbrot() string {
	const (
		xmin, ymin, xmax, ymax = -2, -2, 2, 2
		width, height          = 1024, 1024
		file                   = "mandel.png"
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrotColored(z))
		}
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.Remove(file)
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Cannot open file %s", file)
		return file
	}

	png.Encode(f, img)
	f.Close()
	log.Printf("Wrote %s", file)

	return file
}
