//go:build embed_frontend
// +build embed_frontend

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/dist
var staticFiles embed.FS

// getStaticFileSystem returns the embedded static file system
func getStaticFileSystem() http.FileSystem {
	// Extract the embedded filesystem
	fsys, err := fs.Sub(staticFiles, "web/dist")
	if err != nil {
		// If web/dist doesn't exist in embed, return empty filesystem
		return http.FS(staticFiles)
	}
	return http.FS(fsys)
}

