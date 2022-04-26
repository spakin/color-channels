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
}

// ParseCommandLine parses the command line into a Parameters struct.  It
// aborts on error.
func ParseCommandLine(p *Parameters) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <image-file>...\n", os.Args[0])
		fmt.Fprint(flag.CommandLine.Output(), "Options:\n\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&p.OutputName, "o", "", "Name of output file (default standard output)")
	flag.StringVar(&p.ColorSpace, "space", "hcl",
		`Color space in which to interpret the input channels ("hcl", "hsl", "hsluv", "luv", "lab", "srgb", or "xyy"`)
	flag.Parse()
	p.InputNames = flag.Args()
}

func main() {
	notify = log.New(os.Stderr, os.Args[0]+": ", 0)
	var p Parameters
	ParseCommandLine(&p)
	MergeChannels(&p)
}
