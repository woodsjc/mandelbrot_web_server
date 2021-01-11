package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"
)

const (
	iterations = 200
)

func mandelbrotGreyscale(z complex128) color.Color {
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

func mandelbrotIterations(z complex128) int {
	var v complex128
	var n int

	for n = 0; n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return n + 1
		}
	}

	return n
}

func antiAliasMB(z complex128, height int, width int) color.Color {
	miny := complex(0, float64(1/(2*height)))
	minx := complex(float64(1/(2*width)), 0)
	total := 0

	var sample [4]complex128
	sample[0] = z - minx - miny
	sample[1] = z - minx + miny
	sample[2] = z + minx - miny
	sample[3] = z + minx + miny

	for _, i := range sample {
		total += mandelbrotIterations(i)
	}

	tmp := int(math.Round(float64(total / len(sample))))
	if tmp == iterations {
		return color.Black
	}
	return palette.Plan9[uint8(tmp)]
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
			//img.Set(px, py, mandelbrotColored(z))
			img.Set(px, py, antiAliasMB(z, height, width))
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
