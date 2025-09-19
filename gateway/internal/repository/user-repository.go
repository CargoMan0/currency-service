package repository

import (
	"context"
	"github.com/BernsteinMondy/currency-service/gateway/internal/models"
	"github.com/BernsteinMondy/currency-service/gateway/internal/repository/errors"
	"sync"
)

type UserRepository struct {
	users map[string]models.RepoUser
	mu    *sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]models.RepoUser),
		mu:    &sync.RWMutex{},
	}
}

func (ur *UserRepository) SaveUser(_ context.Context, user models.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	_, exists := ur.users[user.Login]
	if exists {
		return errors.ErrRepoAlreadyExists
	}

	ur.users[user.Login] = models.RepoUser{
		Login:    user.Login,
		Password: user.Password,
	}

	return nil
}

func (ur *UserRepository) GetUserByLogin(_ context.Context, login string) (models.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	user, exists := ur.users[login]
	if !exists {
		return models.User{}, errors.ErrRepoNotFound
	}

	return models.User{
		Login:    user.Login,
		Password: user.Password,
	}, nil
}
