package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func serveExport(w http.ResponseWriter, r *http.Request, name string) {
	if !strings.HasSuffix(name, ".zip") {
		writeStatus(w, http.StatusBadRequest)
		return
	}
	baseDir := filepath.Join(os.TempDir(), "line-sticker-exports")
	path := filepath.Join(baseDir, name)
	if _, err := os.Stat(path); err != nil {
		writeStatus(w, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, path)
}
