# xpm

Pure Go XPM image encoder/decoder.

## Install
```sh
go get -u github.com/coalaura/xpm
```

## Usage

### Decode
```go
package main

import (
	"fmt"
	"image"
	"os"

	"github.com/coalaura/xpm"
)

func main() {
	f, _ := os.Open("test.xpm")
	defer f.Close()

	img, err := xpm.Decode(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(img.Bounds())
}
```

### Encode
```go
package main

import (
	"image"
	"image/color"
	"os"

	"github.com/coalaura/xpm"
)

func main() {
	img := image.NewGray(image.Rect(0, 0, 8, 8))
	img.SetGray(3, 3, color.Gray{Y: 255})

	f, _ := os.Create("out.xpm")
	defer f.Close()

	xpm.Encode(f, img, xpm.XPMOptions{Name: "test"})
}
```