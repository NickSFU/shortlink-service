package shortlink

import (
	"context"

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

// Save сохраняет ссылку
func (r *Repository) Save(code, url string) error {
	query := `
		INSERT INTO short_links (code, original_url)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(context.Background(), query, code, url)
	return err
}

// GetByCode получает ссылку полностью
func (r *Repository) GetByCode(code string) (*ShortLink, error) {
	query := `
		SELECT id, code, original_url
		FROM short_links
		WHERE code = $1
		AND is_deleted IS false
	`

	var link ShortLink

	err := r.db.QueryRow(context.Background(), query, code).Scan(
		&link.ID,
		&link.Code,
		&link.OriginalURL,
	)
	if err != nil {
		return nil, err
	}

	return &link, nil
}
