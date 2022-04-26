// This file provides functions for reading and writing images.

package main

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"

	_ "github.com/spakin/netpbm"
)

// ReadImage reads an arbitrary image from a named file.  It aborts on error.
func ReadImage(fn string) image.Image {
	// Read the input image.
	r, err := os.Open(fn)
	if err != nil {
		notify.Fatal(err)
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	if err != nil {
		notify.Fatal(err)
	}
	return img
}

// ReadGrayscaleImage reads a grayscale image from a named file.  It aborts on
// error.
func ReadGrayscaleImage(fn string) *image.Gray {
	// Read a generic image.
	img := ReadImage(fn)

	// Convert the image to grayscale.
	bnds := img.Bounds()
	gray := image.NewGray(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}

// WritePNG writes an arbitrary image to a named PNG file.  If the file is "",
// write to standard output.
func WritePNG(fn string, img image.Image) error {
	var w io.Writer = os.Stdout
	if fn != "" {
		f, err := os.Create(fn)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	err := png.Encode(w, img)
	if err != nil {
		return err
	}
	return nil
}
