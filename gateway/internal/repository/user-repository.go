package repository

import (
	"context"
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"sync"
)

type User struct {
	Login    string
	Password string
}

type UserRepository struct {
	users map[string]User
	mu    *sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]User),
		mu:    &sync.RWMutex{},
	}
}

func (ur *UserRepository) SaveUser(_ context.Context, user service.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	_, exists := ur.users[user.Login]
	if exists {
		return ErrRepoAlreadyExists
	}

	ur.users[user.Login] = User(user)
	return nil
}

func (ur *UserRepository) GetUserByLogin(_ context.Context, login string) (service.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	user, exists := ur.users[login]
	if !exists {
		return service.User{}, ErrRepoNotFound
	}

	return service.User(user), nil
}
