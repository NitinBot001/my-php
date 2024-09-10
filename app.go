package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// VideoInfo holds the information extracted by yt-dlp
type VideoInfo struct {
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	AudioURL     string `json:"audio_url"`
}

// extractVideoInfo extracts the video title, thumbnail, and audio link using yt-dlp
func extractVideoInfo(videoURL string) (*VideoInfo, error) {
	// Run yt-dlp command to extract JSON metadata
	cmd := exec.Command("yt-dlp", "-j", "--no-playlist", videoURL)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	// Extract needed info
	title := result["title"].(string)
	thumbnail := result["thumbnail"].(string)

	// Extract audio URL using yt-dlp best audio format
	cmdAudio := exec.Command("yt-dlp", "-f", "bestaudio", "-g", videoURL)
	audioURL, err := cmdAudio.Output()
	if err != nil {
		return nil, err
	}

	// Prepare video information struct
	videoInfo := &VideoInfo{
		Title:        title,
		ThumbnailURL: thumbnail,
		AudioURL:     strings.TrimSpace(string(audioURL)),
	}

	return videoInfo, nil
}

// handleExtract is the HTTP handler function to process YouTube URLs
func handleExtract(w http.ResponseWriter, r *http.Request) {
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	// Extract video information
	videoInfo, err := extractVideoInfo(videoURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract video info: %v", err), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videoInfo)
}

func main() {
	http.HandleFunc("/extract", handleExtract)
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
