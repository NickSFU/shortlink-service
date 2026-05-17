package shortlink

import (
	"context"
	"encoding/json"
	stdErrors "errors"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"

	"github.com/NickSFU/shortlink-service/internal/apperror"
)

// бизнес-логика
type Service struct {
	repo  *Repository
	cache *redis.Client
}

// конструктор
func NewService(
	repo *Repository,
	cache *redis.Client,
) *Service {

	return &Service{
		repo:  repo,
		cache: cache,
	}
}

// создаёт короткую ссылку
func (s *Service) CreateShortLink(
	userID int,
	url string,
) (string, error) {

	if url == "" {
		return "",
			apperror.ErrInvalidURL
	}

	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {

		code := generateCode(6)

		err := s.repo.Save(
			userID,
			code,
			url,
		)

		if err == nil {
			return code, nil
		}

		// duplicate code
		var pgErr *pgconn.PgError

		if stdErrors.As(err, &pgErr) {

			if pgErr.Code == "23505" {
				continue
			}
		}

		return "", err
	}

	return "",
		stdErrors.New(
			"failed to generate unique code",
		)
}

// получает ссылку
func (s *Service) GetLink(
	code string,
) (*ShortLink, error) {

	ctx := context.Background()

	// Redis cache
	cachedData, err := s.cache.Get(
		ctx,
		code,
	).Result()

	if err == nil {

		var link ShortLink

		err = json.Unmarshal(
			[]byte(cachedData),
			&link,
		)

		if err == nil {
			return &link, nil
		}
	}

	// PostgreSQL
	link, err := s.repo.GetByCode(code)

	if err != nil {

		return nil,
			apperror.ErrLinkNotFound
	}

	// cache set
	data, err := json.Marshal(link)

	if err == nil {

		_ = s.cache.Set(
			ctx,
			code,
			data,
			time.Hour,
		).Err()
	}

	return link, nil
}

// список ссылок пользователя
func (s *Service) GetUserLinks(
	userID int,
) ([]ShortLink, error) {

	return s.repo.GetByUserID(userID)
}

// soft delete
func (s *Service) DeleteLink(
	code string,
) error {

	err := s.repo.SoftDelete(code)

	if err != nil {
		return err
	}

	// cache cleanup
	_ = s.cache.Del(
		context.Background(),
		code,
	).Err()

	return nil
}

// update original url
func (s *Service) UpdateLink(
	code string,
	newURL string,
) error {

	if newURL == "" {

		return apperror.ErrInvalidURL
	}

	err := s.repo.UpdateURL(
		code,
		newURL,
	)

	if err != nil {
		return err
	}

	// remove stale cache
	_ = s.cache.Del(
		context.Background(),
		code,
	).Err()

	return nil
}

// генерация кода
func generateCode(length int) string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(
		rand.NewSource(
			time.Now().UnixNano(),
		),
	)

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}
