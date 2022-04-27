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
	"strings"
	"unicode"
)

// notify is used to output error messages.
var notify *log.Logger

// Parameters encapsulates all program parameters.
type Parameters struct {
	InputNames     []string // Input file names
	OutputName     string   // Output file names
	OrigColorSpace string   // Color-space name as written by the user
	ColorSpace     string   // Color-space name
	Split          bool     // true: split; false: merge
	Alpha          bool     // true: split/merge an alpha layer: false: don't
}

// colorSpaceList is a list of acceptable color spaces, represented as
// lowercase strings.
var colorSpaceList = []string{
	"cmyk",
	"hcl",
	"hsl",
	"hsluv",
	"lab",
	"linrgb",
	"luv",
	"rgb",
	"srgb",
	"xyy",
	"ycbcr",
}

// colorSpaceString is a list of acceptable color spaces, represented as a
// single, lowercase string with "or" before the final color-space name.
var colorSpaceString string

// init initializes colorSpaceString from colorSpaceList
func init() {
	quoted := make([]string, len(colorSpaceList))
	for i, cs := range colorSpaceList {
		quoted[i] = `"` + cs + `"`
	}
	ncs := len(quoted)
	quoted[ncs-1] = "or " + quoted[ncs-1] // Assume at least 3 color spaces.
	colorSpaceString = strings.Join(quoted, ", ")
	colorSpaceString += `, with an optional "a" suffix`
}

// cleanColorSpaceName maps a color-space name to lowercase and removes spaces
// and asterisks.  Hence, "L*a*b*" maps to "lab", for example.
func cleanColorSpaceName(cs string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, cs)
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
	flag.StringVar(&p.OrigColorSpace, "space", "rgb",
		"Color space in which to interpret the input channels ("+colorSpaceString+")")
	split := flag.Bool("split", false, "Split a color image into one grayscale image per color channel")
	merge := flag.Bool("merge", false, "Merge one grayscale image per color channel into a single color image")
	flag.Parse()
	p.InputNames = flag.Args()

	// Validate the use of the --split and --merge arguments.
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

	// Ensure a valid color space was designated.  Determine if an alpha
	// channel should be used.
	p.ColorSpace = cleanColorSpaceName(p.OrigColorSpace)
	var validCS bool
	for _, cs := range colorSpaceList {
		if p.ColorSpace == cs {
			validCS = true
			break
		}
	}
	if !validCS && len(p.ColorSpace) >= 1 && p.ColorSpace[len(p.ColorSpace)-1] == 'a' {
		// Second chance: Look for an alpha channel.
		opaque := p.ColorSpace[:len(p.ColorSpace)-1]
		for _, cs := range colorSpaceList {
			if opaque == cs {
				validCS = true
				p.ColorSpace = opaque
				p.Alpha = true
				break
			}
		}
	}
	if !validCS {
		notify.Fatalf("--space requires one of %s (not %q)",
			colorSpaceString, p.OrigColorSpace)
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
