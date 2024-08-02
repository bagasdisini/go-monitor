package main

import (
	"embed"
	"fmt"
	"go-monitor/internal/monitor"
	"go-monitor/internal/websocket"
	"net/http"
	"time"
)

//go:embed index.html
var content embed.FS

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", websocket.HandleConnections)
	http.HandleFunc("/api/monitor/", HandleURLSubmission)

	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	data, err := content.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Could not read index.html", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func HandleURLSubmission(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url != "" {
		go monitor.StartMonitoring(500*time.Millisecond, url)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
