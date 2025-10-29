package main

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	listenAddr    = ":8080"
	storageDir    = "./data/images"
	publicPrefix  = "/files/"
	maxUploadSize = 100
)

type ImageItem struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
}

func main() {
	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/images", withCORS(imagesHandler))
	mux.HandleFunc(publicPrefix, withCORS(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(publicPrefix, http.FileServer(http.Dir(storageDir))).ServeHTTP(w, r)
	}))

	log.Printf("Starting server on %s...", listenAddr)

	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleList(w, r)
	case http.MethodPost:
		handleUpload(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(storageDir)
	if err != nil {
		http.Error(w, "Failed to read storage directory", http.StatusInternalServerError)
		return
	}

	var images []ImageItem
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		low := strings.ToLower(name)
		if !(strings.HasSuffix(low, ".png")) {
			continue
		}
		fi, err := os.Stat(filepath.Join(storageDir, low))
		if err != nil {
			continue
		}

		images = append(images, ImageItem{
			ID:        name,
			URL:       publicPrefix + name,
			Size:      fi.Size(),
			CreatedAt: fi.ModTime().Format(time.RFC3339),
		})
	}

	writeJSON(w, images)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxUploadSize<<20))
	if err := r.ParseMultipartForm(int64(maxUploadSize << 20)); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !isAllowed(header) {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	ext := pickExt(header)
	id := uuid.New().String() + ext
	dstPath := filepath.Join(storageDir, id)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	resp := ImageItem{
		ID:        id,
		URL:       publicPrefix + id,
		Size:      header.Size,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	writeJSON(w, resp)
}

func pickExt(header *multipart.FileHeader) string {
	ct := strings.ToLower(header.Header.Get("Content-Type"))
	switch {
	case strings.Contains(ct, "png"):
		return ".png"
	default:
		return ".bin"
	}
}

func isAllowed(header *multipart.FileHeader) bool {
	ct := header.Header.Get("Content-Type")
	ct = strings.ToLower(ct)
	return strings.HasPrefix(ct, "image/png") || strings.HasPrefix(ct, "image/jpeg") ||
		strings.HasPrefix(ct, "image/webp")
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
