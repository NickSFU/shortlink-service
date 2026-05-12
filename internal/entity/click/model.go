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
	TotalClicks int
	UniqueIPs   int
	LastClickAt *time.Time `json:"last_click_at"`
}
