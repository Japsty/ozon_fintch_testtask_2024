package services

import (
	"Ozon_testtask/internal/models"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrHashErr = errors.New("hashing error: ")

type UserService struct {
	UserRepo models.UserRepo
}

func NewUserService(userRepo models.UserRepo) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (us *UserService) AddUser(ctx context.Context, username, password string) (models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.Join(ErrHashErr, err)
	}
	createdUser, err := us.UserRepo.CreateUser(ctx, uuid.NewString(), username, passwordHash)
	if err != nil {
		return models.User{}, err
	}
	return createdUser, nil
}

func (us *UserService) GetUser(ctx context.Context, username, password string) (models.User, error) {
	gotUser, err := us.UserRepo.GetUser(ctx, username, password)
	if err != nil {
		return models.User{}, err
	}
	return gotUser, nil
}
