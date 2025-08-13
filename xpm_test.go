package xpm

import (
	"crypto/sha256"
	"encoding/hex"
	"image/png"
	"os"
	"testing"
)

const (
	expectedWidth  = 32
	expectedHeight = 32
	expectedHash   = "bea74db55d4d94c3fec1394d20348f99e65fb4cc23bed84fe94ca566a2cbc16d"
)

func Test_DecodeConfig(t *testing.T) {
	file, err := os.OpenFile("image.xpm", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	info, err := DecodeConfig(file)
	if err != nil {
		t.Fatal(err)
	}

	if info.Width != expectedWidth {
		t.Fatalf("expected width of %dpx got: %d", expectedWidth, info.Width)
	}

	if info.Height != expectedHeight {
		t.Fatalf("expected height of %dpx got: %d", expectedHeight, info.Height)
	}
}

func Test_Decode(t *testing.T) {
	file, err := os.OpenFile("image.xpm", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	img, err := Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	bounds := img.Bounds()

	if bounds.Dx() != expectedWidth {
		t.Fatalf("expected width of %dpx got: %d", expectedWidth, bounds.Dx())
	}

	if bounds.Dy() != expectedHeight {
		t.Fatalf("expected height of %dpx got: %d", expectedHeight, bounds.Dy())
	}
}

func Test_Encode(t *testing.T) {
	file, err := os.OpenFile("image.png", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()

	err = Encode(hash, img, XPMOptions{
		Name: "cats",
	})
	if err != nil {
		t.Fatal(err)
	}

	hexHash := hex.EncodeToString(hash.Sum(nil))
	if hexHash != expectedHash {
		t.Fatalf("hash mismatch: %s != %s", hexHash, expectedHash)
	}
}
