package repository

import (
	"context"
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

func (ur *UserRepository) SaveUser(ctx context.Context, user User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	_, exists := ur.users[user.Login]
	if exists {
		return ErrRepoAlreadyExists
	}

	ur.users[user.Login] = user
	return nil
}

func (ur *UserRepository) GetUserByLogin(ctx context.Context, login string) (User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	user, exists := ur.users[login]
	if !exists {
		return User{}, ErrRepoNotFound
	}

	return user, nil
}
