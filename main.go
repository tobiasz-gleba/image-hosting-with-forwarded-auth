package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	authServerBaseURL := os.Getenv("AUTH_SERVER_BASE_URL")
	if authServerBaseURL == "" {
		authServerBaseURL = "http://localhost:8081/auth" // default fallback
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authForwardingHandler(w, r, authServerBaseURL)
	})

	log.Println("Server running on http://localhost")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func authForwardingHandler(w http.ResponseWriter, r *http.Request, authServerBaseURL string) {
	// Build the full auth URL: base + original path and query
	authURL, err := url.Parse(authServerBaseURL)
	if err != nil {
		http.Error(w, "Invalid auth server URL", http.StatusInternalServerError)
		return
	}
	authURL.Path += r.URL.Path
	authURL.RawQuery = r.URL.RawQuery

	// Create the forwarded request
	req, err := http.NewRequest("GET", authURL.String(), nil)
	if err != nil {
		http.Error(w, "Failed to create auth request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Auth server error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Serve image from the request path
	imageVariants := []string{".png", ".jpg", ".jpeg"}
	basePath := strings.TrimPrefix(r.URL.Path, "/")
	for _, ext := range imageVariants {
		fileName := basePath + ext
		file, err := os.Open(fileName)
		if err == nil {
			defer file.Close()

			switch strings.ToLower(filepath.Ext(fileName)) {
			case ".png":
				w.Header().Set("Content-Type", "image/png")
			case ".jpg", ".jpeg":
				w.Header().Set("Content-Type", "image/jpeg")
			}

			io.Copy(w, file)
			return
		}
	}

	http.Error(w, "Image not found", http.StatusInternalServerError)
}
