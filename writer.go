package xpm

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

type XPMOptions struct {
	Name string
}

func Encode(w io.Writer, m image.Image, opts ...XPMOptions) error {
	var name string

	if len(opts) > 0 {
		name = opts[0].Name
	}

	if name == "" {
		name = "image"
	}

	b := m.Bounds()
	width, height := b.Dx(), b.Dy()

	colorKeys := make(map[color.Color]string)

	var palette []color.Color

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := m.At(x, y)

			if _, ok := colorKeys[c]; !ok {
				colorKeys[c] = ""

				palette = append(palette, c)
			}
		}
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

func generateSymbols(n int) []string {
	chars := []rune(" .oO+@#$%&*=-;:>,<1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if n > len(chars) {
		panic("too many colors for simple symbol set")
	}

	syms := make([]string, n)

	for i := 0; i < n; i++ {
		syms[i] = string(chars[i])
	}

	return syms
}
