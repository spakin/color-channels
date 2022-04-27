/*
color-channels splits an image into separate color channels and merges color
channels into a new image.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// notify is used to output error messages.
var notify *log.Logger

// Parameters encapsulates all program parameters.
type Parameters struct {
	InputNames []string // Input file names
	OutputName string   // Output file names
	ColorSpace string   // Color space name
	Split      bool     // true: split; false: merge
}

// ParseCommandLine parses the command line into a Parameters struct.  It
// aborts on error.
func ParseCommandLine(p *Parameters) {
	// Parse the command line.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--merge | --split] [other_options] <image-file>...\n", os.Args[0])
		fmt.Fprint(flag.CommandLine.Output(), "Options:\n\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&p.OutputName, "o", "",
		`Name of output file for --merge (default standard output) or output-file template containing "%s" for --split (no default)`)
	flag.StringVar(&p.ColorSpace, "space", "rgb",
		`Color space in which to interpret the input channels ("hcl", "hsl", "hsluv", "luv", "lab", "linrgb", "rgb", or "xyy"`)
	split := flag.Bool("split", false, "Split a color image into one grayscale image per color channel")
	merge := flag.Bool("merge", false, "Merge one grayscale image per color channel into a single color image")
	flag.Parse()
	p.InputNames = flag.Args()

	// Validate the given arguments.
	switch {
	case *split && *merge:
		notify.Fatal("--split and --merge are mutually exclusive")
	case *split:
		p.Split = true
	case *merge:
		p.Split = false
	case !*split && !*merge:
		notify.Fatal("Exactly one of --split and --merge must be specified")
	}
}

func main() {
	notify = log.New(os.Stderr, os.Args[0]+": ", 0)
	var p Parameters
	ParseCommandLine(&p)
	if p.Split {
		SplitImage(&p)
	} else {
		MergeChannels(&p)
	}
}
