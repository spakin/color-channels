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
	"strconv"
	"strings"
	"unicode"

	"github.com/lucasb-eyer/go-colorful"
)

// notify is used to output error messages.
var notify *log.Logger

// Parameters encapsulates all program parameters.
type Parameters struct {
	InputNames     []string   // Input file names
	OutputName     string     // Output file names
	OrigColorSpace string     // Color-space name as written by the user
	ColorSpace     string     // Color-space name
	Split          bool       // true: split; false: merge
	Alpha          bool       // true: split/merge an alpha layer: false: don't
	WhitePoint     [3]float64 // White reference point as an XYZ color
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
	"xyz",
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

// cleanColorSpaceName maps a color-space name to lowercase and removes
// non-letters.  Hence, "L*a*b*" maps to "lab" and "Y'CbCr" maps to "ycbcr",
// for example.
func cleanColorSpaceName(cs string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, cs)
}

// parseWhitePoint parses a pair of CIE chromaticity coordinates into an XYZ
// color.  It aborts on error.
func parseWhitePoint(s string) [3]float64 {
	// Handle the cases go-colorful supports.
	wp := strings.ToUpper(strings.TrimSpace(s))
	if wp == "D65" {
		return colorful.D65
	}
	if wp == "D50" {
		return colorful.D50
	}

	// Parse the strings into a pair of floating-point numbers.
	toks := strings.FieldsFunc(wp, func(c rune) bool {
		if unicode.IsDigit(c) {
			return false
		}
		switch c {
		case '+', '-', 'E', '.':
			return false
		default:
			return true
		}
	})
	if len(toks) != 2 {
		notify.Fatalf(`Failed to parse %q as either "D65", "D50", or a pair of floating-point numbers`, s)
	}
	x, err := strconv.ParseFloat(toks[0], 64)
	if err != nil || x < 0.0 || x > 1.0 {
		notify.Fatalf("Failed to parse %q as a floating-point number in [0.0, 1.0]", toks[0])
	}
	y, err := strconv.ParseFloat(toks[1], 64)
	if err != nil || y <= 0.0 || y > 1.0 {
		notify.Fatalf("Failed to parse %q as a floating-point number in (0.0, 1.0]", toks[0])
	}
	if x+y > 1.0 {
		notify.Fatalf("%s + %s must be less than or equal to 1.0", toks[0], toks[1])
	}

	// Convert from (x, y) to XYZ.
	z := 1.0 - x - y
	return [3]float64{x / y, 1.0, z / y}
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
	white := flag.String("white", "D65",
		`White-point CIE chromaticity coordinates (two numbers in [0.0, 1.0]) or "D65" or "D50", used for hcl, lab, and luv`)
	flag.Parse()
	p.InputNames = flag.Args()
	p.WhitePoint = parseWhitePoint(*white)

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
