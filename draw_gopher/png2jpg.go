package main

import (
	"flag"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

// png2jpg: convert given PNG image to JPEG image.
//
// Example:
// go run png2jpg.go -out gopher.jpg gopher.png

func main() {
	outFilename := flag.String("out", "", "output file name")
	quality := flag.Int("q", 80, "output JPEG quality")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("expected exactly 1 input file name")
	}

	filename := flag.Args()[0] // Image to convert

	switch {
	case *quality < 0:
		*quality = 0
	case *quality > 100:
		*quality = 100
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Panicf("open input image: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Panicf("decode input image: %v", err)
	}

	if *outFilename == "" {
		*outFilename = strings.ReplaceAll(filename, "png", "jpg")
	}

	outFile, err := os.Create(*outFilename)
	if err != nil {
		log.Panicf("create file: %v", err)
	}
	defer outFile.Close()

	opts := &jpeg.Options{Quality: *quality}
	if err := jpeg.Encode(outFile, img, opts); err != nil {
		log.Panicf("encode: %v", err)
	}
}
