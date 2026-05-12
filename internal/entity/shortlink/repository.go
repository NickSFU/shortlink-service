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
func (r *Repository) Save(
	userID int,
	code string,
	url string,
) error {
	query := `
		INSERT INTO short_links (
			user_id,
			code,
			original_url
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		userID,
		code,
		url,
	)

	return err
}

// GetByCode получает ссылку полностью
func (r *Repository) GetByCode(code string) (*ShortLink, error) {
	query := `
		SELECT id, user_id, code, original_url
		FROM short_links
		WHERE code = $1
		AND is_deleted IS false
	`

	var link ShortLink

	err := r.db.QueryRow(context.Background(), query, code).Scan(
		&link.ID,
		&link.UserID,
		&link.Code,
		&link.OriginalURL,
	)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *Repository) GetByUserID(
	userID int,
) ([]ShortLink, error) {
	query := `
		SELECT
			id,
			user_id,
			code,
			original_url
		FROM short_links
		WHERE user_id = $1
		AND is_deleted = false
		ORDER BY id DESC
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []ShortLink

	for rows.Next() {
		var link ShortLink

		err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.Code,
			&link.OriginalURL,
		)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

func (r *Repository) SoftDelete(
	code string,
) error {
	query := `
		UPDATE short_links
		SET is_deleted = true
		WHERE code = $1
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		code,
	)

	return err
}

func (r *Repository) UpdateURL(
	code string,
	newURL string,
) error {
	query := `
		UPDATE short_links
		SET original_url = $1
		WHERE code = $2
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		newURL,
		code,
	)

	return err
}
