package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/currency-service/gateway/internal/models"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user models.User) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type AuthClient interface {
	GenerateToken(ctx context.Context, login string) (string, error)
	ValidateToken(ctx context.Context, token string) error
}

type AuthService struct {
	repository UserRepository
	authClient AuthClient
}

func NewAuthService(repository UserRepository, authClient AuthClient) *AuthService {
	return &AuthService{
		repository: repository,
		authClient: authClient,
	}
}

func (s *AuthService) Register(ctx context.Context, login, password string) error {
	user := models.User{
		Login:    login,
		Password: password,
	}

	err := s.repository.SaveUser(ctx, user)
	if err != nil {
		if errors.Is(err, ErrRepoAlreadyExists) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("repository: save user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.repository.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, ErrRepoNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("repository: get user by login: %w", err)
	}

	if user.Password != password {
		return "", ErrInvalidCredentials
	}

	token, err := s.authClient.GenerateToken(ctx, login)
	if err != nil {
		return "", fmt.Errorf("repository: generate token: %w", err)
	}

	return token, nil
}
