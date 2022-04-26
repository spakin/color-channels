// This file provides routines for merging color channels to produce a new
// image.

package main

import (
	"image"

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

// MergeSRGB merges R, G, and B channels into a single image.
func MergeSRGB(imgs []*image.Gray) image.Image {
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
