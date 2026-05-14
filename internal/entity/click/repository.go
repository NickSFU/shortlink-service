package click

import (
	"context"

	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(click *Click) error {
	query := `
		INSERT INTO clicks (
			short_link_id,
			ip,
			user_agent,
			referer
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		click.ShortLinkID,
		click.IP,
		click.UserAgent,
		click.Referer,
	)

	return err
}

// CountByLink считает количество переходов
func (r *Repository) CountByLink(shortLinkID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM clicks
		WHERE short_link_id = $1
	`

	var count int

	err := r.db.QueryRow(
		context.Background(),
		query,
		shortLinkID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetStats(
	linkID int,
) (*LinkStats, error) {

	query := `
		SELECT
			COUNT(*) as total_clicks,
			COUNT(DISTINCT ip) as unique_ips,
			MAX(created_at) as last_click_at
		FROM clicks
		WHERE short_link_id = $1
	`

	var stats LinkStats
	var lastClick sql.NullTime

	err := r.db.QueryRow(
		context.Background(),
		query,
		linkID,
	).Scan(
		&stats.TotalClicks,
		&stats.UniqueIPs,
		&lastClick,
	)

	if err != nil {
		return nil, err
	}

	if lastClick.Valid {
		stats.LastClickAt = &lastClick.Time
	}

	// peak hour
	peakHourQuery := `
		SELECT
			EXTRACT(HOUR FROM created_at) as hour
		FROM clicks
		WHERE short_link_id = $1
		GROUP BY hour
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`

	_ = r.db.QueryRow(
		context.Background(),
		peakHourQuery,
		linkID,
	).Scan(&stats.PeakHour)

	// top referer
	topRefererQuery := `
		SELECT referer
		FROM clicks
		WHERE short_link_id = $1
		AND referer != ''
		GROUP BY referer
		ORDER BY COUNT(*) DESC
		LIMIT 1
	`

	_ = r.db.QueryRow(
		context.Background(),
		topRefererQuery,
		linkID,
	).Scan(&stats.TopReferer)

	return &stats, nil
}
