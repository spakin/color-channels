// This file provides functions for splitting an image into separate channels.

package main

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// A ImageInfo represents a channel name and image data.
type ImageInfo struct {
	Name  string      // Channel name
	Image *image.Gray // Grayscale image representing a channel
}

// toGrayVal converts a float64 in [0.0, 1.0] to a color.Gray, clamping if
// necessary.
func toGrayVal(f float64) color.Gray {
	if f < 0.0 {
		return color.Gray{Y: 0}
	}
	if f > 1.0 {
		return color.Gray{Y: 255}
	}
	return color.Gray{Y: uint8(f * 255.0)}
}

// allocGrays allocates an array of N grayscale images of a given size.
func allocGrays(bnds image.Rectangle, n int) []*image.Gray {
	grays := make([]*image.Gray, n)
	for i := range grays {
		grays[i] = image.NewGray(bnds)
	}
	return grays
}

// splitAny is a helper function for the various Split* functions.  It performs
// all the boilerplate code, invoking a color space-specific function for each
// pixel.
func splitAny(img image.Image, names []string,
	fn func(colorful.Color) []float64) []ImageInfo {
	bnds := img.Bounds()
	grays := allocGrays(bnds, len(names))
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			clr, _ := colorful.MakeColor(img.At(x, y))
			for i, f := range fn(clr) {
				grays[i].Set(x, y, toGrayVal(f))
			}
		}
	}
	result := make([]ImageInfo, len(names))
	for i, nm := range names {
		result[i].Name = nm
		result[i].Image = grays[i]
	}
	return result
}

// SplitHCL splits a color image into separate H, C, and L channels.
func SplitHCL(img image.Image) []ImageInfo {
	return splitAny(img, []string{"H", "C", "L"},
		func(clr colorful.Color) []float64 {
			h, c, l := clr.Hcl()
			return []float64{h / 360.0, c, l}
		})
}

// SplitLab splits a color image into separate L*, a*, and b* channels.
func SplitLab(img image.Image) []ImageInfo {
	return splitAny(img, []string{"L", "a", "b"},
		func(clr colorful.Color) []float64 {
			l, a, b := clr.Lab()
			return []float64{l, (a + 1.0) / 2.0, (b + 1.0) / 2.0}
		})
}

// SplitLuv splits a color image into separate L*, u*, and v* channels.
func SplitLuv(img image.Image) []ImageInfo {
	return splitAny(img, []string{"L", "u", "v"},
		func(clr colorful.Color) []float64 {
			l, u, v := clr.Luv()
			return []float64{l, (u + 1.0) / 2.0, (v + 1.0) / 2.0}
		})
}

// SplitImage splits an image into separate channel images.  It aborts on error.
func SplitImage(p *Parameters) {
	// Ensure we have exactly one input file.
	if len(p.InputNames) != 1 {
		notify.Fatalf("Expected 1 input file but saw %d", len(p.InputNames))
	}

	// Ensure the output file contains a "%s".
	if p.OutputName == "" {
		notify.Fatal("An output-file template must be specified when --split is used")
	}
	if !strings.Contains(p.OutputName, "%s") {
		notify.Fatalf(`With --split, the output file must contain "%%s"`)
	}

	// Read the input image.
	inImg := ReadImage(p.InputNames[0])

	// Split the input image into multiple grayscale images.
	var outImgs []ImageInfo
	switch p.ColorSpace {
	case "hcl":
		outImgs = SplitHCL(inImg)
	case "lab":
		outImgs = SplitLab(inImg)
	case "luv":
		outImgs = SplitLuv(inImg)
	default:
		notify.Fatal("Invalid argument to --space")
	}

	// Write each channel to a separate grayscale file.
	for _, info := range outImgs {
		name := fmt.Sprintf(p.OutputName, info.Name)
		WritePNG(name, info.Image)
	}
}
