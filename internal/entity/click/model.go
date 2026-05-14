package click

import "time"

type Click struct {
	ID          int
	ShortLinkID int
	IP          string
	UserAgent   string
	Referer     string
	CreatedAt   time.Time
}

type LinkStats struct {
	TotalClicks int        `json:"total_clicks"`
	UniqueIPs   int        `json:"unique_ips"`
	LastClickAt *time.Time `json:"last_click_at"`

	PeakHour   int    `json:"peak_hour"`
	TopReferer string `json:"top_referer"`
}
