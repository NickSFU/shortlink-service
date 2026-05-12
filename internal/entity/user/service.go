package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/NickSFU/shortlink-service/internal/auth"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(
	email string,
	password string,
) error {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.repo.Create(
		email,
		string(hash),
	)
}

func (s *Service) Login(
	email string,
	password string,
) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateToken(user.ID)
}
