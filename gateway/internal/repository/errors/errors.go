package errors

import "errors"

var (
	ErrRepoAlreadyExists = errors.New("repo: already exists")
	ErrRepoNotFound      = errors.New("repo: not found")
)
