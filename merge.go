// This file provides routines for merging color channels to produce a new
// image.

package main

import (
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// MergeHCL merges H, C, and L channels into a single image.
func MergeHCL(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			h := float64(imgs[0].GrayAt(x, y).Y) * 360.0 / 255.0
			c := float64(imgs[1].GrayAt(x, y).Y) / 255.0
			l := float64(imgs[2].GrayAt(x, y).Y) / 255.0
			clr := colorful.Hcl(h, c, l).Clamped()
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeLab merges L*, a*, and b* channels into a single image.
func MergeLab(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			L := float64(imgs[0].GrayAt(x, y).Y) / 255.0
			a := float64(imgs[1].GrayAt(x, y).Y)*2.0/255.0 - 1.0
			b := float64(imgs[2].GrayAt(x, y).Y)*2.0/255.0 - 1.0
			clr := colorful.Lab(L, a, b).Clamped()
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeLuv merges L*, u*, and v* channels into a single image.
func MergeLuv(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			L := float64(imgs[0].GrayAt(x, y).Y) / 255.0
			u := float64(imgs[1].GrayAt(x, y).Y)*2.0/255.0 - 1.0
			v := float64(imgs[2].GrayAt(x, y).Y)*2.0/255.0 - 1.0
			clr := colorful.Luv(L, u, v).Clamped()
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeXyy merges x, y, and Y channels into a single image.
func MergeXyy(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for r := bnds.Min.Y; r < bnds.Max.Y; r++ {
		for c := bnds.Min.X; c < bnds.Max.X; c++ {
			x := float64(imgs[0].GrayAt(c, r).Y) / 255.0
			y := float64(imgs[1].GrayAt(c, r).Y) / 255.0
			Y := float64(imgs[2].GrayAt(c, r).Y) / 255.0
			clr := colorful.Xyy(x, y, Y).Clamped()
			merged.Set(c, r, clr)
		}
	}
	return merged
}

// MergeHSL merges H, S, and L channels into a single image.
func MergeHSL(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			h := float64(imgs[0].GrayAt(x, y).Y) * 360.0 / 255.0
			s := float64(imgs[1].GrayAt(x, y).Y) / 255.0
			l := float64(imgs[2].GrayAt(x, y).Y) / 255.0
			clr := colorful.Hsl(h, s, l).Clamped()
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeHSLuv merges H, S, and L channels into a single image.
func MergeHSLuv(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			h := float64(imgs[0].GrayAt(x, y).Y) * 360.0 / 255.0
			s := float64(imgs[1].GrayAt(x, y).Y) / 255.0
			l := float64(imgs[2].GrayAt(x, y).Y) / 255.0
			clr := colorful.HSLuv(h, s, l).Clamped()
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeLinRGB merges R, G, and B channels into a single image.
func MergeLinRGB(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			r := float64(imgs[0].GrayAt(x, y).Y) / 255.0
			g := float64(imgs[1].GrayAt(x, y).Y) / 255.0
			b := float64(imgs[2].GrayAt(x, y).Y) / 255.0
			clr := colorful.LinearRgb(r, g, b).Clamped()
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeRGB merges R, G, and B channels into a single image.
func MergeRGB(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			r := imgs[0].GrayAt(x, y).Y
			g := imgs[1].GrayAt(x, y).Y
			b := imgs[2].GrayAt(x, y).Y
			clr := color.NRGBA{r, g, b, 255}
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeSRGB merges R, G, and B channels into a single image.
func MergeSRGB(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			r := float64(imgs[0].GrayAt(x, y).Y) / 255.0
			g := float64(imgs[1].GrayAt(x, y).Y) / 255.0
			b := float64(imgs[2].GrayAt(x, y).Y) / 255.0
			clr := colorful.Color{R: r, G: g, B: b}
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeCMYK merges C, M, Y, and K channels into a single image.
func MergeCMYK(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			c := imgs[0].GrayAt(x, y).Y
			m := imgs[1].GrayAt(x, y).Y
			w := imgs[2].GrayAt(x, y).Y // y is already taken.
			k := imgs[3].GrayAt(x, y).Y
			r, g, b := color.CMYKToRGB(c, m, w, k)
			clr := color.NRGBA{r, g, b, 255}
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// MergeYCbCr merges Y, Cb, and Cr channels into a single image.
func MergeYCbCr(imgs []*image.Gray) image.Image {
	bnds := imgs[0].Bounds()
	merged := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			l := imgs[0].GrayAt(x, y).Y // y is already taken.
			cb := imgs[1].GrayAt(x, y).Y
			cr := imgs[2].GrayAt(x, y).Y
			r, g, b := color.YCbCrToRGB(l, cb, cr)
			clr := color.NRGBA{r, g, b, 255}
			merged.Set(x, y, clr)
		}
	}
	return merged
}

// AddAlpha replaces an image's alpha channel with a separately specified alpha
// channel.
func AddAlpha(img image.Image, alpha *image.Gray) image.Image {
	bnds := img.Bounds()
	newImg := image.NewNRGBA(bnds)
	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			clr := img.At(x, y)
			nrgba := color.NRGBAModel.Convert(clr).(color.NRGBA)
			nrgba.A = alpha.GrayAt(x, y).Y
			newImg.Set(x, y, nrgba)
		}
	}
	return newImg
}

// MergeChannels merges the input files into a single output file.  It aborts
// on error.
func MergeChannels(p *Parameters) {
	// Ensure we have the correct number of input files.
	nIn := len(p.InputNames)
	wrongArgsFmt := "Expected %d input files for --space=%q but saw %d"
	numAlpha := 0
	if p.Alpha {
		numAlpha = 1
	}
	switch p.ColorSpace {
	case "cmyk":
		if nIn != 4+numAlpha {
			notify.Fatalf(wrongArgsFmt, 4+numAlpha, p.OrigColorSpace, nIn)
		}
	default:
		if nIn != 3+numAlpha {
			notify.Fatalf(wrongArgsFmt, 3+numAlpha, p.OrigColorSpace, nIn)
		}
	}

	// Read all the color-channel images.
	channels := make([]*image.Gray, 0, 4)
	for _, fn := range p.InputNames {
		g := ReadGrayscaleImage(fn)
		channels = append(channels, g)
	}

	// Ensure that all channels have the same bounds.
	bnds := channels[0].Bounds()
	for _, g := range channels {
		if g.Bounds() != bnds {
			notify.Fatal("All input images must have the same dimensions")
		}
	}

	// Merge the channels.
	var merged image.Image
	switch p.ColorSpace {
	case "cmyk":
		merged = MergeCMYK(channels)
	case "hcl":
		merged = MergeHCL(channels)
	case "hsl":
		merged = MergeHSL(channels)
	case "hsluv":
		merged = MergeHSLuv(channels)
	case "lab":
		merged = MergeLab(channels)
	case "linrgb":
		merged = MergeLinRGB(channels)
	case "luv":
		merged = MergeLuv(channels)
	case "rgb":
		merged = MergeRGB(channels)
	case "srgb":
		merged = MergeSRGB(channels)
	case "xyy":
		merged = MergeXyy(channels)
	case "ycbcr":
		merged = MergeYCbCr(channels)
	default:
		panic("Internal error: unimplemented color space")
	}

	// If an alpha channel was included, insert it into the image.
	if p.Alpha {
		merged = AddAlpha(merged, channels[len(channels)-1])
	}

	// Write the result to a file.
	err := WritePNG(p.OutputName, merged)
	if err != nil {
		notify.Fatal(err)
	}
}
