package inmem

import (
	"Ozon_testtask/internal/models"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

var (
	ErrNoUser     = errors.New("user not found")
	ErrBadPass    = errors.New("invalid password")
	ErrUserExists = errors.New("username already exists")
)

type UserMemoryRepository struct {
	data map[string]models.User
	mu   *sync.RWMutex
}

func NewUserRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: map[string]models.User{},
		mu:   &sync.RWMutex{},
	}
}

func (ur *UserMemoryRepository) CreateUser(_ context.Context, userID, login string, passHash []byte) (models.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	for _, user := range ur.data {
		if user.Username == login {
			return models.User{}, ErrUserExists
		}
	}

	u := models.User{
		Username:     login,
		PasswordHash: passHash,
		ID:           userID,
	}

	ur.data[u.ID] = u

	return u, nil
}

func (ur *UserMemoryRepository) GetUser(_ context.Context, login, pass string) (models.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	var foundUser *models.User
	for _, user := range ur.data {
		if user.Username == login {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		return models.User{}, ErrNoUser
	}

	err := bcrypt.CompareHashAndPassword(foundUser.PasswordHash, []byte(pass))
	if err != nil {
		return models.User{}, ErrBadPass
	}

	return *foundUser, nil
}
