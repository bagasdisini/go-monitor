package models

import "time"

type WebsiteStatus struct {
	URL     string        `json:"url"`
	Status  string        `json:"status"`
	Latency time.Duration `json:"latency"`
}
