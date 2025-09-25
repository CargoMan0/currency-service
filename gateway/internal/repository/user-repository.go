package repository

import (
	"context"
	"github.com/BernsteinMondy/currency-service/gateway/internal/repository/errors"
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"sync"
)

type UserRepository struct {
	users map[string]repoUser
	mu    *sync.RWMutex
}

var _ service.UserRepository = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]repoUser),
		mu:    &sync.RWMutex{},
	}
}

func (ur *UserRepository) SaveUser(_ context.Context, user service.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	_, exists := ur.users[user.Login]
	if exists {
		return errors.ErrRepoAlreadyExists
	}

	ur.users[user.Login] = repoUser{
		Login:    user.Login,
		Password: user.Password,
	}

	return nil
}

func (ur *UserRepository) GetUserByLogin(_ context.Context, login string) (service.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	user, exists := ur.users[login]
	if !exists {
		return service.User{}, errors.ErrRepoNotFound
	}

	return service.User{
		Login:    user.Login,
		Password: user.Password,
	}, nil
}
