package service

import (
	"context"
	"github.com/BernsteinMondy/currency-service/gateway/internal/repository"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user repository.User) error
	GetUser() (repository.User, error)
}

type AuthService struct {
	repository repository.UserRepository
}

func NewAuthService(repository repository.UserRepository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}
