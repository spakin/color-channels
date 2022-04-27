// This file provides functions for splitting an image into separate channels.

package main

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"sync"

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
	var wg sync.WaitGroup
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		// Concurrently process all rows
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := bnds.Min.X; x < bnds.Max.X; x++ {
				clr, _ := colorful.MakeColor(img.At(x, y))
				for i, f := range fn(clr) {
					grays[i].Set(x, y, toGrayVal(f))
				}
			}
		}(y)
	}
	wg.Wait()
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

// SplitXyy splits a color image into separate x, y, and Y channels.
func SplitXyy(img image.Image) []ImageInfo {
	return splitAny(img, []string{"x", "y", "Y"},
		func(clr colorful.Color) []float64 {
			x, y, Y := clr.Xyy()
			return []float64{x, y, Y}
		})
}

// SplitHSL splits a color image into separate H, S, and L channels.
func SplitHSL(img image.Image) []ImageInfo {
	return splitAny(img, []string{"H", "S", "L"},
		func(clr colorful.Color) []float64 {
			h, s, l := clr.Hsl()
			return []float64{h / 360.0, s, l}
		})
}

// SplitHSLuv splits a color image into separate H, S, and L channels.
func SplitHSLuv(img image.Image) []ImageInfo {
	return splitAny(img, []string{"H", "S", "L"},
		func(clr colorful.Color) []float64 {
			h, s, l := clr.HSLuv()
			return []float64{h / 360.0, s, l}
		})
}

// SplitLinRGB splits a color image into separate R, G, and B channels.
func SplitLinRGB(img image.Image) []ImageInfo {
	return splitAny(img, []string{"R", "G", "B"},
		func(clr colorful.Color) []float64 {
			r, g, b := clr.LinearRgb()
			return []float64{r, g, b}
		})
}

// SplitRGB splits a color image into separate R, G, and B channels.
func SplitRGB(img image.Image) []ImageInfo {
	return splitAny(img, []string{"R", "G", "B"},
		func(clr colorful.Color) []float64 {
			ri, gi, bi := clr.RGB255()
			r := float64(ri) / 255.0
			g := float64(gi) / 255.0
			b := float64(bi) / 255.0
			return []float64{r, g, b}
		})
}

// SplitSRGB splits a color image into separate R, G, and B channels.
func SplitSRGB(img image.Image) []ImageInfo {
	return splitAny(img, []string{"R", "G", "B"},
		func(clr colorful.Color) []float64 {
			return []float64{clr.R, clr.G, clr.B}
		})
}

// SplitCMYK splits a color image into separate C, M, Y, and K channels.
func SplitCMYK(img image.Image) []ImageInfo {
	return splitAny(img, []string{"C", "M", "Y", "K"},
		func(clr colorful.Color) []float64 {
			ri, gi, bi := clr.RGB255()
			ci, mi, yi, ki := color.RGBToCMYK(ri, gi, bi)
			c := float64(ci) / 255.0
			m := float64(mi) / 255.0
			y := float64(yi) / 255.0
			k := float64(ki) / 255.0
			return []float64{c, m, y, k}
		})
}

// SplitYCbCr splits a color image into separate Y, Cb, and Cr channels.
func SplitYCbCr(img image.Image) []ImageInfo {
	return splitAny(img, []string{"Y", "Cb", "Cr"},
		func(clr colorful.Color) []float64 {
			ri, gi, bi := clr.RGB255()
			yi, cbi, cri := color.RGBToYCbCr(ri, gi, bi)
			l := float64(yi) / 255.0 // y is already taken.
			cb := float64(cbi) / 255.0
			cr := float64(cri) / 255.0
			return []float64{l, cb, cr}
		})
}

// ExtractAlpha extracts an image's alpha channel and returns it as an
// ImageInfo.
func ExtractAlpha(img image.Image) ImageInfo {
	bnds := img.Bounds()
	gray := image.NewGray(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			clr := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			gray.SetGray(x, y, color.Gray{Y: clr.A})
		}
	}
	return ImageInfo{
		Name:  "alpha",
		Image: gray,
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
	case "cmyk":
		outImgs = SplitCMYK(inImg)
	case "hcl":
		outImgs = SplitHCL(inImg)
	case "hsl":
		outImgs = SplitHSL(inImg)
	case "hsluv":
		outImgs = SplitHSLuv(inImg)
	case "lab":
		outImgs = SplitLab(inImg)
	case "linrgb":
		outImgs = SplitLinRGB(inImg)
	case "luv":
		outImgs = SplitLuv(inImg)
	case "rgb":
		outImgs = SplitRGB(inImg)
	case "srgb":
		outImgs = SplitSRGB(inImg)
	case "xyy":
		outImgs = SplitXyy(inImg)
	case "ycbcr":
		outImgs = SplitYCbCr(inImg)
	default:
		panic("Internal error: unimplemented color space")
	}

	// Optionally include an alpha channel.
	if p.Alpha {
		outImgs = append(outImgs, ExtractAlpha(inImg))
	}

	// Write each channel to a separate grayscale file.
	for _, info := range outImgs {
		name := fmt.Sprintf(p.OutputName, info.Name)
		WritePNG(name, info.Image)
	}
}
