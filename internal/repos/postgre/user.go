package postgre

import (
	"Ozon_testtask/internal/models"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	ErrNoUser     = errors.New("user not found")
	ErrBadPass    = errors.New("invalid password")
	ErrUserExists = errors.New("username already exists")
)

type UserMySQLRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserMySQLRepository {
	return &UserMySQLRepository{db: db}
}

func (repo *UserMySQLRepository) CreateUser(ctx context.Context, userID, login string, password []byte) (models.User, error) {
	var exists bool
	err := repo.db.QueryRow(querries.CheckUserExists, login).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists = false
		}
		return models.User{}, err
	}

	if exists {
		return models.User{}, ErrUserExists
	}

	var createdUser models.User
	var dbPasswordHash string
	_, err = repo.db.Exec(querries.CreateUser, userID, login, string(password))
	if err != nil {
		return models.User{}, err
	}

	err = repo.db.QueryRow(querries.GetUserByID, userID).Scan(&createdUser.ID, &createdUser.Username, &dbPasswordHash)
	if err != nil {
		return models.User{}, err
	}

	createdUser.PasswordHash = []byte(dbPasswordHash)

	return createdUser, nil
}

func (repo *UserMySQLRepository) GetUser(ctx context.Context, login, pass string) (models.User, error) {
	var exists bool
	err := repo.db.QueryRow(querries.CheckUserExists, login).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists = false
		}
		return models.User{}, err
	}

	if !exists {
		return models.User{}, ErrNoUser
	}

	var foundUser models.DBUser
	err = repo.db.QueryRow(querries.GetUserByLogin, login).Scan(&foundUser.Username, &foundUser.PasswordHash, &foundUser.ID)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:     foundUser.Username,
		PasswordHash: []byte(foundUser.PasswordHash),
		ID:           foundUser.ID,
	}
	log.Print(user)

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pass))
	if err != nil {
		return models.User{}, ErrBadPass
	}

	return user, err
}
