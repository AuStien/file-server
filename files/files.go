package files

import (
	"path/filepath"
	"strings"
)

var VaildFiletypes = []string{"jpeg", "jpg", "png", "gif", "mp4", "mp3", "wav", "avi", "mkv", "mov", "webm"}

type File struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func toVaildFiles(entries []string) []string {
	validEntries := []string{}

	for _, e := range entries {
		if isValidFiletype(filepath.Ext(e)) {
			validEntries = append(validEntries, e)
		}
	}

	return validEntries
}

func isValidFiletype(ext string) bool {
	ext = strings.TrimPrefix(ext, ".")

	for _, f := range VaildFiletypes {
		if f == ext {
			return true
		}
	}

	return false
}
