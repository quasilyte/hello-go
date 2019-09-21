package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
)

// invert: invert PNG image colors.
//
// Example:
// go run invert.go gopher.png

func main() {
	outFilename := flag.String("out", "inverted.png", "output file name")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("expected exactly 1 input file name")
	}

	filename := flag.Args()[0] // Image to invert

	f, err := os.Open(filename)
	if err != nil {
		log.Panicf("open input image: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Panicf("decode input image: %v", err)
	}

	// Type-assert ("cast") generic image to NRGBA to access individual pixels.
	dst, ok := img.(*image.NRGBA)
	if !ok {
		log.Panicf("only NRGBA images are supported")
	}

	bounds := dst.Bounds()
	for y := 0; y < bounds.Size().Y; y++ {
		i := y * dst.Stride
		for x := 0; x < bounds.Size().X; x++ {
			// Use i+4 to capture alpha as well.
			// Since we leave it "as is", we take a slice of 3
			// components: RGB. Alpha would be in d[3].
			d := dst.Pix[i : i+3 : i+3]

			// Invert colors.
			d[0] = 255 - d[0] // R
			d[1] = 255 - d[1] // G
			d[2] = 255 - d[2] // B

			i += 4
		}
	}

	outFile, err := os.Create(*outFilename)
	if err != nil {
		log.Panicf("create file: %v", err)
	}
	defer outFile.Close()
	if err := png.Encode(outFile, img); err != nil {
		log.Panicf("encode: %v", err)
	}
}
