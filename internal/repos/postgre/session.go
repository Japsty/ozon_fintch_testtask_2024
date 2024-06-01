package postgre

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre/querries"
	"database/sql"
	"errors"
)

var ErrTokenNotFound = errors.New("token not found")

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (repo *SessionRepository) Create(timeToExpire int64, tokenString, userID, username string) error {
	_, err := repo.DB.Exec(
		querries.CreateSession,
		userID,
		username,
		timeToExpire,
		tokenString,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SessionRepository) GetSession(token string) (*model.Session, error) {
	var session model.Session

	err := repo.DB.QueryRow(
		querries.GetSessionByToken,
		token,
	).Scan(&session.ID, &session.UserID, &session.Username, &session.Exp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (repo *SessionRepository) Delete(token string) error {
	_, err := repo.DB.Exec(
		querries.DeleteSession,
		token,
	)
	return err
}
