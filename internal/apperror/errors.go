package apperror

import "errors"

var (
	ErrUserExists = errors.New(
		"user already exists",
	)

	ErrInvalidCredentials = errors.New(
		"invalid credentials",
	)

	ErrLinkNotFound = errors.New(
		"link not found",
	)

	ErrForbidden = errors.New(
		"forbidden",
	)

	ErrInvalidURL = errors.New(
		"invalid url",
	)
)
