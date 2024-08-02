package monitor

import (
	"go-monitor/internal/websocket"
	"go-monitor/pkg/models"
	"net/http"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

func normalizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return url
}

func checkWebsite(url string) models.WebsiteStatus {
	url = normalizeURL(url)
	start := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(start)

	status := models.WebsiteStatus{
		URL:     url,
		Latency: time.Duration(latency.Milliseconds()),
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		status.Status = "DOWN"
	} else {
		status.Status = "UP"
	}

	return status
}

func StartMonitoring(interval time.Duration, url string) {
	ticker := time.NewTicker(interval)
	start := time.Now()
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			status := checkWebsite(url)
			websocket.BroadcastStatus(status)
			mu.Unlock()
		}

		if time.Since(start) > 1*time.Hour {
			ticker.Stop()
			break
		}
	}
}
