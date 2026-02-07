package api

import (
	"bytes"
	"image"
	"image/png"
)

func validatePNGSize(data []byte, w, h int) bool {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return false
	}
	b := img.Bounds()
	return b.Dx() == w && b.Dy() == h
}
