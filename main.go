package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func authForwardingHandler(w http.ResponseWriter, r *http.Request, authServerBaseURL string) {
	// Build the full auth URL: base + original path and query
	authURL, err := url.Parse(authServerBaseURL)
	if err != nil {
		http.Error(w, "Invalid auth server URL", http.StatusInternalServerError)
		return
	}
	authURL.Path += r.URL.Path
	authURL.RawQuery = r.URL.RawQuery

	// // Create the forwarded request
	req, err := http.NewRequest("GET", authURL.String(), nil)
	if err != nil {
		http.Error(w, "Failed to create auth request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Auth server error for request %s: %v", r.URL.Path, err)
		http.Error(w, "Auth server error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read response body for logging purposes
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unauthorized access attempt for image: %s", r.URL.Path)
		log.Printf("Auth server response: Status=%d, Body=%s", resp.StatusCode, string(respBody))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Serve image from the request path
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "/app/static" // default fallback
	}
	basePath := strings.TrimPrefix(r.URL.Path, "/")
	fileName := filepath.Join(staticDir, basePath)
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

	log.Printf("Image not found: %s", fileName)
	http.Error(w, "Image not found", http.StatusInternalServerError)
}

func main() {
	authServerBaseURL := os.Getenv("AUTH_SERVER_BASE_URL")
	if authServerBaseURL == "" {
		authServerBaseURL = "https://agentix.pl" // default fallback
	}
	log.Printf("Using auth server base URL: %s", authServerBaseURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authForwardingHandler(w, r, authServerBaseURL)
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
