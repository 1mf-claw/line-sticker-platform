package api

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"io"
	"net/http"
	"time"

	"golang.org/x/image/draw"
)

const (
	stickerWidth  = 370
	stickerHeight = 320
)

func normalizeStickerImage(imageURL string) (string, error) {
	return normalizeImageToSize(imageURL, stickerWidth, stickerHeight)
}

func normalizeImageToSize(imageURL string, targetW, targetH int) (string, error) {
	if imageURL == "" {
		return "", errors.New("empty image url")
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", errors.New("image fetch failed")
	}
	limited := io.LimitReader(resp.Body, 10<<20)
	src, _, err := image.Decode(limited)
	if err != nil {
		return "", err
	}
	bw := src.Bounds().Dx()
	bh := src.Bounds().Dy()
	if bw == 0 || bh == 0 {
		return "", errors.New("invalid image")
	}

	scaleW := float64(targetW) / float64(bw)
	scaleH := float64(targetH) / float64(bh)
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}
	newW := int(float64(bw) * scale)
	newH := int(float64(bh) * scale)
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}
	offX := (targetW - newW) / 2
	offY := (targetH - newH) / 2

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	dstRect := image.Rect(offX, offY, offX+newW, offY+newH)
	draw.CatmullRom.Scale(dst, dstRect, src, src.Bounds(), draw.Over, nil)

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, dst); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + encoded, nil
}
