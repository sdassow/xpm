package xpm

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"
)

func DecodeConfig(r io.Reader) (image.Config, error) {
	lines, err := readXPMStringsFast(r)
	if err != nil {
		return image.Config{}, err
	}

	if len(lines) == 0 {
		return image.Config{}, fmt.Errorf("xpm: empty data")
	}

	parts := strings.Fields(lines[0])
	if len(parts) < 4 {
		return image.Config{}, fmt.Errorf("xpm: invalid header")
	}

	w, _ := strconv.Atoi(parts[0])
	h, _ := strconv.Atoi(parts[1])

	return image.Config{
		ColorModel: color.RGBAModel,
		Width:      w,
		Height:     h,
	}, nil
}

func Decode(r io.Reader) (image.Image, error) {
	lines, err := readXPMStringsFast(r)
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("xpm: empty data")
	}

	parts := strings.Fields(lines[0])
	if len(parts) < 4 {
		return nil, fmt.Errorf("xpm: invalid header")
	}

	w, _ := strconv.Atoi(parts[0])
	h, _ := strconv.Atoi(parts[1])

	numColors, _ := strconv.Atoi(parts[2])
	charsPerPixel, _ := strconv.Atoi(parts[3])

	colorMap := make(map[string]color.Color, numColors)

	for i := 0; i < numColors; i++ {
		line := lines[1+i]

		key := line[:charsPerPixel]

		fields := strings.Fields(line[charsPerPixel:])
		if len(fields) < 2 || fields[0] != "c" {
			return nil, fmt.Errorf("xpm: invalid color line: %q", line)
		}

		colorMap[key] = parseColor(fields[1])
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		line := lines[1+numColors+y]

		for x := 0; x < w; x++ {
			key := line[x*charsPerPixel : (x+1)*charsPerPixel]

			c, ok := colorMap[key]
			if !ok {
				return nil, fmt.Errorf("xpm: unknown color key %q", key)
			}

			img.Set(x, y, c)
		}
	}

	return img, nil
}

func parseColor(s string) color.Color {
	if s == "None" {
		return color.RGBA{0, 0, 0, 0}
	}

	if strings.HasPrefix(s, "#") {
		hex := s[1:]

		if len(hex) == 6 {
			r, _ := strconv.ParseUint(hex[0:2], 16, 8)
			g, _ := strconv.ParseUint(hex[2:4], 16, 8)
			b, _ := strconv.ParseUint(hex[4:6], 16, 8)

			return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		}
	}

	return color.Black
}

func readXPMStringsFast(r io.Reader) ([]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 0, len(data))

	var (
		inQuote   bool
		inComment bool
	)

	for i := 0; i < len(data); i++ {
		b := data[i]

		if !inQuote && i+1 < len(data) && data[i] == '/' && data[i+1] == '*' {
			inComment = true

			i++

			continue
		}

		if inComment && i+1 < len(data) && data[i] == '*' && data[i+1] == '/' {
			inComment = false

			i++

			continue
		}

		if inComment {
			continue
		}

		if b == '"' {
			inQuote = !inQuote

			continue
		}

		if inQuote {
			out = append(out, b)
		} else if b == '\n' {
			out = append(out, '\n')
		}
	}

	rawLines := bytes.Split(out, []byte{'\n'})

	lines := make([]string, 0, len(rawLines))

	for _, l := range rawLines {
		if len(l) > 0 {
			lines = append(lines, string(l))
		}
	}

	return lines, nil
}
