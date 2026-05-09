package shortlink

import (
	"errors"
	"math/rand"
	"time"

	"github.com/jackc/pgconn"
)

// бизнес-логи
type Service struct {
	repo *Repository
}

// конструктор
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// создаёт короткую ссылку
func (s *Service) CreateShortLink(url string) (string, error) {
	if url == "" {
		return "", errors.New("url is empty")
	}

	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		code := generateCode(6)

		err := s.repo.Save(code, url)
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
	return s.repo.GetByCode(code)
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
