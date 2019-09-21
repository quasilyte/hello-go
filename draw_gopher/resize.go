package main

import (
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

// resize: resize given image to the specified size.
// If either of dimensions (-w or -h) are 0, aspect ratios are preserved.
//
// Example:
// go run resize.go -w 100 gopher.png

func main() {
	// This program uses almost the same steps as compose.go.

	// Flags binding and decoding.
	// Using Uint instead of Int because "Resize" expects unsigned values.
	width := flag.Uint("w", 640, "output image width in pixels")
	height := flag.Uint("h", 0, "output image height in pixels")
	outFilename := flag.String("out", "resized.png", "output file name")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("expected exactly 1 file name")
	}
	filename := flag.Args()[0] // Image to resize

	f, err := os.Open(filename)
	if err != nil {
		log.Panicf("open input image: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Panicf("decode input image: %v", err)
	}

	// Resizing itself.
	resizedImg := resize.Resize(*width, *height, img, resize.Bicubic)

	outFile, err := os.Create(*outFilename)
	if err != nil {
		log.Panicf("create file: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, resizedImg); err != nil {
		log.Panicf("encode: %v", err)
	}
}
