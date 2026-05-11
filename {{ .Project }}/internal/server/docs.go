package server

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// DocsHandler serves static files from docs directory.
func DocsHandler() http.Handler {
	candidates := make([]string, 0, 3)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		candidates = append(candidates, "docs")
	} else {
		candidates = append(candidates, filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "docs")))
	}

	if exePath, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exePath), "docs"))
	}

	candidates = append(candidates, "docs")

	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return http.StripPrefix("/docs/", http.FileServer(http.Dir(dir)))
		}
	}

	return http.StripPrefix("/docs/", http.FileServer(http.Dir("docs")))
}
