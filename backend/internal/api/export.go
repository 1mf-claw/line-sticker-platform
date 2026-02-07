package api

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func buildExportZip(projectID string, stickers []Sticker) (string, error) {
	if projectID == "" {
		return "", errors.New("missing project id")
	}
	if len(stickers) == 0 {
		return "", errors.New("no stickers")
	}

	baseDir := filepath.Join(os.TempDir(), "line-sticker-exports")
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}

	zipPath := filepath.Join(baseDir, fmt.Sprintf("%s.zip", projectID))
	f, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for i, s := range stickers {
		url := s.TransparentURL
		if url == "" {
			url = s.ImageURL
		}
		data, err := fetchPNG(url)
		if err != nil {
			continue
		}
		name := fmt.Sprintf("%02d.png", i+1)
		w, err := zw.Create(name)
		if err != nil {
			continue
		}
		_, _ = w.Write(data)
	}
	if err := zw.Close(); err != nil {
		return "", err
	}
	return zipPath, nil
}

func fetchPNG(url string) ([]byte, error) {
	if url == "" {
		return nil, errors.New("empty url")
	}
	if len(url) > 30 && url[:22] == "data:image/png;base64," {
		// data url
		b64 := url[len("data:image/png;base64,"):]
		return decodeBase64(b64)
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, errors.New("fetch failed")
	}
	limited := io.LimitReader(resp.Body, 10<<20)
	img, _, err := image.Decode(limited)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeBase64(v string) ([]byte, error) {
	buf := make([]byte, len(v))
	_, err := base64.StdEncoding.Decode(buf, []byte(v))
	if err != nil {
		return nil, err
	}
	return buf, nil
}
