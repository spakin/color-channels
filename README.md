Color Channels
==============

[![Go Report Card](https://goreportcard.com/badge/github.com/spakin/color-channels)](https://goreportcard.com/report/github.com/spakin/color-channels)

Color Channels is a simple, command-line program that can

1. separate a color image into multiple grayscale images, each representing a single color channel, and
2. combine multiple grayscale images, each representing a single color channel, into a unified color image.

Features
--------

Color Channels supports the following color spaces:

* [CMYK](https://en.wikipedia.org/wiki/CMYK_color_model)
* [HCL](https://en.wikipedia.org/wiki/HCL_color_space)
* [HSL](https://en.wikipedia.org/wiki/HSL_and_HSV)
* [HSLuv](https://en.wikipedia.org/wiki/HSLuv)
* [L\*a\*b\*](https://en.wikipedia.org/wiki/CIELAB_color_space)
* [Linear RGB](https://www.sjbrown.co.uk/posts/gamma-correct-rendering/)
* [L\*u\*v\*](https://en.wikipedia.org/wiki/CIELUV)
* [RGB](https://en.wikipedia.org/wiki/RGB_color_spaces)
* [sRGB](https://en.wikipedia.org/wiki/SRGB)
* [xyY](https://en.wikipedia.org/wiki/CIE_1931_color_space)
* [XYZ](https://en.wikipedia.org/wiki/CIE_1931_color_space)
* [Y'CbCr](https://en.wikipedia.org/wiki/YCbCr)

An alpha (opacity) channel is supported for *all* of the above.

The user can specify an explicit [white point](https://en.wikipedia.org/wiki/White_point) for HCL, L\*a\*b\*, and L\*u\*v\* conversions.

The program accepts images provided in [PNG](https://en.wikipedia.org/wiki/Portable_Network_Graphics), [JPEG](https://en.wikipedia.org/wiki/JPEG), [GIF](https://en.wikipedia.org/wiki/GIF), or any of the [Netpbm](https://en.wikipedia.org/wiki/Netpbm) formats.

Unrepresentable colors are clamped gracefully to representable colors.

Installation
------------

Color Channels is written in the [Go programming language](https://golang.org/) so you will need to install Go to compile it.

Once Go is installed, Color Channels can be downloaded, built, and installed into `$GOPATH/bin/` simply by running
```bash
go install github.com/spakin/color-channels@latest
```
An alternative install location can be specified by first setting the `GOBIN` environment variable (e.g., `export GOBIN=/usr/local/bin`).

Usage
-----

### Basic operation

Run `color-channels --help` for a usage summary.  In short, one of `--split` or `--merge` must be specified.  For example,
```bash
color-channels --split --space=HCL -o channel-%s.png input-image.jpg
```
reads `input-image.jpg` and generates `channel-H.png`, representing the hue channel, `channel-C.png`, representing the chroma channel, and `channel-L.png`, representing the luminance channel.  Output images always are written in PNG regardless of the input-image's format.

The channel images from the preceding command can be recombined (typically after transforming them in some manner) using `--merge`:
```bash
color-channels --merge --space=HCL -o output-image.png channel-H.png channel-C.png channel-L.png
```

### A more concrete example

Here's a sample image, courtesy of https://file-examples.com/:

![input-image](https://user-images.githubusercontent.com/650041/165878104-330de79f-48a0-42ad-959d-27e257a92f30.jpg)

We can split the above into luma (gamma-corrected luminance), blue-difference, and red-difference channels with
```bash
color-channels --split --space="Y'CbCr" -o channel-%s.png example.jpg
```
![channel-Y](https://user-images.githubusercontent.com/650041/165878503-27f4afa6-e03c-4fe2-bfd0-b89feeeec299.jpg)
![channel-Cb](https://user-images.githubusercontent.com/650041/165878469-91cfcd46-f8ab-46bf-976e-877c1f9bd024.jpg)
![channel-Cr](https://user-images.githubusercontent.com/650041/165878470-1e8a8ba8-af2d-4579-8e99-e3d717e11033.jpg)

Let's swap the blue-difference and red-difference channels when recombining the above to see what happens:
```bash
color-channels --merge --space="Y'CbCr" -o output-image.png channel-Y.png channel-Cr.png channel-Cb.png 
```
![output-image](https://user-images.githubusercontent.com/650041/165878472-f69c9f3d-d410-4399-9050-70ac9416a149.jpg)

### Color spaces

The `--space` option allows color-space names to include arbitrary punctuation and capitalization.  That is, `--space="L*a*b*"` and `--space=lab` are treated identically.

Appending an `A` to any color-space name includes an alpha channel (named `alpha` on output).

### Advanced usage

As a more advanced example, a white point can be specified via its *x* and *y* chromaticity coordinates.  `color-channels` currently honors the white point for only a few color spaces, though.  The default white point is [D65](https://en.wikipedia.org/wiki/Illuminant_D65).  As an example, the [F2 standard illuminant](https://en.wikipedia.org/wiki/Standard_illuminant) (cool white fluorescent) with a 2° standard observer can be requested with `--white="0.37208 0.37529"`.

Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott+clrch@pakin.org*
