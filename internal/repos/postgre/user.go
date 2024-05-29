package postgre

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	ErrNoUser     = errors.New("user not found")
	ErrBadPass    = errors.New("invalid password")
	ErrUserExists = errors.New("username already exists")
)

type UserMySQLRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserMySQLRepository {
	return &UserMySQLRepository{db: db}
}

func (repo *UserMySQLRepository) CreateUser(ctx context.Context, userID, login string, password []byte) (model.User, error) {
	var exists bool
	err := repo.db.QueryRow(ctx, querries.CheckUserExists, login).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists = false
		}
		return model.User{}, err
	}

	if exists {
		return model.User{}, ErrUserExists
	}

	var createdUser model.User
	var dbPasswordHash string
	_, err = repo.db.Exec(ctx, querries.CreateUser, userID, login, string(password))
	if err != nil {
		return model.User{}, err
	}

	err = repo.db.QueryRow(ctx, querries.GetUserByID, userID).Scan(&createdUser.ID, &createdUser.Username, &dbPasswordHash)
	if err != nil {
		return model.User{}, err
	}

	createdUser.PasswordHash = []byte(dbPasswordHash)

	return createdUser, nil
}

func (repo *UserMySQLRepository) GetUser(ctx context.Context, login, pass string) (model.User, error) {
	var exists bool
	err := repo.db.QueryRow(ctx, querries.CheckUserExists, login).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists = false
		}
		return model.User{}, err
	}

	if !exists {
		return model.User{}, ErrNoUser
	}

	var foundUser model.DBUser
	err = repo.db.QueryRow(ctx, querries.GetUserByLogin, login).Scan(&foundUser.Username, &foundUser.PasswordHash, &foundUser.ID)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Username:     foundUser.Username,
		PasswordHash: []byte(foundUser.PasswordHash),
		ID:           foundUser.ID,
	}
	log.Print(user)

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pass))
	if err != nil {
		return model.User{}, ErrBadPass
	}

	return user, err
}
