package atlas

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Handler(urlPathPrefix string, rootDir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path, found := strings.CutPrefix(r.URL.Path, urlPathPrefix)
		if !found {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		ServeMedia(w, r, filepath.Join(rootDir, path))
	}
}

func ServeMedia(w http.ResponseWriter, r *http.Request, path string) {
	fileInfo, err := os.Stat(path)
	if err != nil || fileInfo.IsDir() {
		http.Error(w, "Media file not found", http.StatusNotFound)
		return
	}

	contentType := getMediaContentType(path)
	if contentType == "" {
		http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "Failed to open media file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")

	rc := http.NewResponseController(w)
	err = rc.EnableFullDuplex()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isVideo := isVideoContentType(contentType)
	if isVideo {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=86400")
	}

	if isVideo {
		fileName := filepath.Base(path)
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileName))

		// Prevent duplication of buffering for large files
		// and disable automatic compression for videos
		if fileInfo.Size() > 10*1024*1024 { // files >10MB
			// Disable potential automatic compression
			w.Header().Set("Content-Encoding", "identity")
		}

		if r.Header.Get("Range") != "" {
			rc.SetWriteDeadline(time.Now().Add(5 * time.Minute))
		}
	} else if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// Apply compression only to non-video files and only for full requests
		if r.Header.Get("Range") == "" {
			w.Header().Set("Content-Encoding", "gzip")
			gzWriter := gzip.NewWriter(w)
			defer gzWriter.Close()

			w.WriteHeader(http.StatusOK)
			io.Copy(gzWriter, file)
			return
		}
	}

	w.Header().Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))
	http.ServeContent(w, r, filepath.Base(path), fileInfo.ModTime(), file)
}

func getMediaContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".bmp":
		return "image/bmp"
	case ".tiff", ".tif":
		return "image/tiff"
	case ".ico":
		return "image/x-icon"
	// Video
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogv":
		return "video/ogg"
	case ".mov":
		return "video/quicktime"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".flv":
		return "video/x-flv"
	case ".m4v":
		return "video/x-m4v"
	case ".ts":
		return "video/mp2t"
	case ".3gp":
		return "video/3gpp"
	default:
		return ""
	}
}

func isVideoContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "video/")
}
