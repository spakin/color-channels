// This file provides functions for splitting an image into separate channels.

package main

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

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

// A ImageInfo represents a channel name and image data.
type ImageInfo struct {
	Name  string      // Channel name
	Image *image.Gray // Grayscale image representing a channel
}

// SplitHCL splits a color image into separate H, C, and L channels.
func SplitHCL(img image.Image) []ImageInfo {
	// Prepare the output images.
	bnds := img.Bounds()
	var grays [3]*image.Gray
	grays[0] = image.NewGray(bnds)
	grays[1] = image.NewGray(bnds)
	grays[2] = image.NewGray(bnds)

	// Convert each pixel in turn.
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			clr, _ := colorful.MakeColor(img.At(x, y))
			h, c, l := clr.Hcl()
			grays[0].Set(x, y, toGrayVal(h/360.0))
			grays[1].Set(x, y, toGrayVal(c))
			grays[2].Set(x, y, toGrayVal(l))
		}
	}

	// Return the color channels.
	return []ImageInfo{
		{Name: "H", Image: grays[0]},
		{Name: "C", Image: grays[1]},
		{Name: "L", Image: grays[2]},
	}
}

// SplitLab splits a color image into separate L*, a*, and b* channels.
func SplitLab(img image.Image) []ImageInfo {
	// Prepare the output images.
	bnds := img.Bounds()
	var grays [3]*image.Gray
	grays[0] = image.NewGray(bnds)
	grays[1] = image.NewGray(bnds)
	grays[2] = image.NewGray(bnds)

	// Convert each pixel in turn.
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			clr, _ := colorful.MakeColor(img.At(x, y))
			l, a, b := clr.Lab()
			grays[0].Set(x, y, toGrayVal(l))
			grays[1].Set(x, y, toGrayVal((a+1.0)/2.0))
			grays[2].Set(x, y, toGrayVal((b+1.0)/2.0))
		}
	}

	// Return the color channels.
	return []ImageInfo{
		{Name: "L", Image: grays[0]},
		{Name: "a", Image: grays[1]},
		{Name: "b", Image: grays[2]},
	}
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
	default:
		notify.Fatal("Invalid argument to --space")
	}

	// Write each channel to a separate grayscale file.
	for _, info := range outImgs {
		name := fmt.Sprintf(p.OutputName, info.Name)
		WritePNG(name, info.Image)
	}
}
