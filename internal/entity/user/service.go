package user

import (
	stdErrors "errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/NickSFU/shortlink-service/internal/apperror"
	"github.com/NickSFU/shortlink-service/internal/auth"
)

type Service struct {
	repo *Repository
}

func NewService(
	repo *Repository,
) *Service {

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

	err = s.repo.Create(
		email,
		string(hash),
	)

	if err != nil {

		var pgErr *pgconn.PgError

		if stdErrors.As(
			err,
			&pgErr,
		) {

			// duplicate email
			if pgErr.Code == "23505" {

				return apperror.ErrUserExists
			}
		}

		return err
	}

	return nil
}

func (s *Service) Login(
	email string,
	password string,
) (string, error) {

	user, err := s.repo.GetByEmail(email)

	if err != nil {

		return "",
			apperror.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {

		return "",
			apperror.ErrInvalidCredentials
	}

	return auth.GenerateToken(user.ID)
}
