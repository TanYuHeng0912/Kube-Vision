//go:build !embed_frontend
// +build !embed_frontend

package main

import (
	"net/http"
	"os"
)

// getStaticFileSystem returns a filesystem-based static file system
// Used when frontend is not embedded (development mode)
// This function is kept for potential future use
//nolint:unused
func getStaticFileSystem() http.FileSystem {
	// Check if web/dist exists
	if _, err := os.Stat("web/dist"); os.IsNotExist(err) {
		// Return empty filesystem if dist doesn't exist
		return http.Dir(".")
	}
	return http.Dir("web/dist")
}



