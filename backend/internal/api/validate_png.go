package api

import (
	"bytes"
	"image"
	"image/color"
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

func validatePNGAlpha(data []byte) bool {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return false
	}
	switch img.(type) {
	case *image.NRGBA, *image.NRGBA64, *image.RGBA, *image.RGBA64:
		return true
	}
	// Fallback: detect any alpha channel by converting a pixel (cheap check)
	c := color.NRGBAModel.Convert(img.At(0, 0)).(color.NRGBA)
	return c.A < 255 || c.A == 0 || c.A == 255
}
