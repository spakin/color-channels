/*
merge-channels merges separately provided channels to produce a new image.
*/
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	_ "github.com/spakin/netpbm"
)

// notify is used to output error messages.
var notify *log.Logger

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

// MergeHCL merges hue, chroma, and lightness channels into a single image.
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

func main() {
	// Read three images from the command line.
	notify = log.New(os.Stderr, os.Args[0]+": ", 0)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <hue.png> <chroma.png> <lightness.png>\n", os.Args[0])
		fmt.Fprint(flag.CommandLine.Output(), "Options:\n\n")
		flag.PrintDefaults()
	}
	outName := flag.String("o", "", "Name of output stereogram file (default standard output)")
	space := flag.String("space", "hcl",
		`Color space in which to interpret the input channels ("hcl", "lab", "luv", or "xyy")`)
	flag.Parse()
	if flag.NArg() < 3 {
		flag.Usage()
		os.Exit(1)
	}

	// Read the three color-channel images.
	channels := make([]*image.Gray, 0, 4)
	for _, fn := range flag.Args() {
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

	// Merge the channels and write the result to a file.
	var merged image.Image
	switch *space {
	case "hcl":
		merged = MergeHCL(channels)
	case "lab":
		merged = MergeLab(channels)
	case "luv":
		merged = MergeLuv(channels)
	case "xyy":
		merged = MergeXyy(channels)
	default:
		notify.Fatal("Invalid argument to --space")
	}
	err := WritePNG(*outName, merged)
	if err != nil {
		notify.Fatal(err)
	}
}
