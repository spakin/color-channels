/*
color-channels splits an image into separate color channels and merges color
channels into a new image.
*/
package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
)

// notify is used to output error messages.
var notify *log.Logger

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
		`Color space in which to interpret the input channels ("hcl", "hsl", "hsluv", "luv", "lab", "srgb", or "xyy"`)
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
	case "hsl":
		merged = MergeHSL(channels)
	case "hsluv":
		merged = MergeHSLuv(channels)
	case "srgb":
		merged = MergeSRGB(channels)
	default:
		notify.Fatal("Invalid argument to --space")
	}
	err := WritePNG(*outName, merged)
	if err != nil {
		notify.Fatal(err)
	}
}
