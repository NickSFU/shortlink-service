package shortlink

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

// бизнес-логи
type Service struct {
	repo  *Repository
	cache *redis.Client
}

// конструктор
func NewService(repo *Repository, cache *redis.Client) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

// создаёт короткую ссылку
func (s *Service) CreateShortLink(userID int, url string) (string, error) {
	if url == "" {
		return "", errors.New("url is empty")
	}

	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		code := generateCode(6)

		err := s.repo.Save(userID, code, url)
		if err == nil {
			return code, nil
		}

		// проверяем duplicate key
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				continue
			}
		}

		return "", err
	}

	return "", errors.New("failed to generate unique code")
}

// возвращает оригинальный URL по коду
//func (s *Service) GetOriginalURL(code string) (string, error) {
//url, err := s.repo.Get(code)
//if err != nil {
//	return "", err
//}

//return url, nil
//}

func (s *Service) GetLink(code string) (*ShortLink, error) {
	ctx := context.Background()

	// пытаемся получить из Redis
	cachedData, err := s.cache.Get(ctx, code).Result()
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

	// если нет в Redis
	link, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// сериализация
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

// генерация случайного кода
func generateCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)

	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}

func (s *Service) GetUserLinks(
	userID int,
) ([]ShortLink, error) {
	return s.repo.GetByUserID(userID)
}

func (s *Service) DeleteLink(
	code string,
) error {
	// удаление из PostgreSQL
	err := s.repo.SoftDelete(code)
	if err != nil {
		return err
	}

	// удаление из Redis
	_ = s.cache.Del(
		context.Background(),
		code,
	).Err()

	return nil
}

func (s *Service) UpdateLink(
	code string,
	newURL string,
) error {
	// обновляем PostgreSQL
	err := s.repo.UpdateURL(
		code,
		newURL,
	)
	if err != nil {
		return err
	}

	// удаляем stale cache
	_ = s.cache.Del(
		context.Background(),
		code,
	).Err()

	return nil
}
