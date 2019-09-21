package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

// Example:
// go run compose.go ears.png body.png eyes.png teeth.png undernose.png nose.png hands.png

func main() {
	// Bind and parse command-line flags.
	// Defaults are enough for gopher-drawing.
	width := flag.Int("w", 490, "output image width in pixels")
	height := flag.Int("h", 600, "output image height in pixels")
	outFilename := flag.String("out", "gopher.png", "output file name")
	flag.Parse()
	filenames := flag.Args() // Gopher parts
	if len(filenames) == 0 {
		log.Fatalf("expected 1 or more command-line arguments")
	}

	// Turn filenames into image objects.
	// Every image is a layer.
	// The order of layers is important.
	var layers []image.Image
	for i, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			log.Panicf("open part[%d]: %v", i, err)
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			log.Panicf("decode part[%d]: %v", i, err)
		}
		layers = append(layers, img)
	}

	// Draw layers, one by one, on a new image (outImage).
	// Note that the first layer is drawn separately.
	bounds := image.Rect(0, 0, *width, *height)
	outImage := image.NewRGBA(bounds)
	draw.Draw(outImage, bounds, layers[0], image.ZP, draw.Src)
	for _, layer := range layers[1:] {
		draw.Draw(outImage, bounds, layer, image.ZP, draw.Over)
	}

	// Write our new image to a file.
	outFile, err := os.Create(*outFilename)
	if err != nil {
		log.Panicf("create file: %v", err)
	}
	defer outFile.Close()
	if err := png.Encode(outFile, outImage); err != nil {
		log.Panicf("encode: %v", err)
	}
}
