package service

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var (
	ErrRepoAlreadyExists = errors.New("repo: already exists")
	ErrRepoNotFound      = errors.New("repo: not found")
)
