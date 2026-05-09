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
