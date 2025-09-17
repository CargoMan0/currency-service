package auth

import "errors"

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTokenGeneration      = errors.New("token generation failed")

	ErrTokenNotFound         = errors.New("token not found in header")
	ErrInvalidOrExpiredToken = errors.New("invalid signature or token expired")
)
