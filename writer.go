package xpm

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"

	"golang.org/x/image/draw"
)

const (
	MaxCharacters     = 80
	AllowedCharacters = " .oO+@#$%&*=-;:>,<1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type XPMOptions struct {
	Name string
}

func Encode(w io.Writer, m image.Image, opts ...XPMOptions) error {
	var name string

	if len(opts) > 0 {
		name = sanitizeName(opts[0].Name)
	}

	if name == "" {
		name = "image"
	}

	b := m.Bounds()
	width, height := b.Dx(), b.Dy()

	palette, colorKeys := generatePalette(m)

	if len(palette) > MaxCharacters {
		m = quantizeImage(m, MaxCharacters)

		palette, colorKeys = generatePalette(m)
	}

	symbols := generateSymbols(len(palette))

	for i, c := range palette {
		colorKeys[c] = symbols[i]
	}

	fmt.Fprintf(w, "/* XPM */\nstatic char * %s[] = {\n", name)
	fmt.Fprintf(w, "\"%d %d %d 1\",\n", width, height, len(palette))

	for _, c := range palette {
		r, g, b, _ := c.RGBA()

		fmt.Fprintf(w, "\"%s c #%02X%02X%02X\",\n", colorKeys[c], uint8(r>>8), uint8(g>>8), uint8(b>>8))
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		var line strings.Builder

		line.WriteByte('"')

		for x := b.Min.X; x < b.Max.X; x++ {
			line.WriteString(colorKeys[m.At(x, y)])
		}

		line.WriteByte('"')

		if y < b.Max.Y-1 {
			line.WriteByte(',')
		}

		line.WriteByte('\n')

		w.Write([]byte(line.String()))
	}

	fmt.Fprintln(w, "};")

	return nil
}

func generatePalette(img image.Image) ([]color.Color, map[color.Color]string) {
	var palette []color.Color

	b := img.Bounds()
	colorKeys := make(map[color.Color]string)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)

			if _, ok := colorKeys[c]; !ok {
				colorKeys[c] = ""

				palette = append(palette, c)
			}
		}
	}

	return palette, colorKeys
}

func generateSymbols(n int) []string {
	chars := []rune(AllowedCharacters)
	if n > len(chars) {
		n = len(chars)
	}

	syms := make([]string, n)

	for i := 0; i < n; i++ {
		syms[i] = string(chars[i])
	}

	return syms
}

func quantizeImage(src image.Image, maxColors int) image.Image {
	b := src.Bounds()

	dst := image.NewPaletted(b, generateSimplePalette(src, maxColors))

	draw.FloydSteinberg.Draw(dst, b, src, image.Point{})

	return dst
}

func generateSimplePalette(img image.Image, maxColors int) color.Palette {
	seen := make(map[color.Color]struct{})

	var palette color.Palette

	b := img.Bounds()

	step := max((b.Dx()*b.Dy())/maxColors, 1)

	count := 0

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if count%step == 0 {
				c := img.At(x, y)

				if _, ok := seen[c]; !ok {
					seen[c] = struct{}{}

					palette = append(palette, c)

					if len(palette) >= maxColors {
						return palette
					}
				}
			}

			count++
		}
	}

	return palette
}

func sanitizeName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "image"
	}

	var out strings.Builder

	for i, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (i > 0 && r >= '0' && r <= '9') {
			out.WriteRune(r)
		} else {
			out.WriteByte('_')
		}
	}

	return out.String()
}
