package users

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	GetUser(userID string) (*User, error)
}

type repository struct {
	*sqlx.DB
}

func (repo *repository) GetUser(userID string) (*User, error) {
	var user User
	if err := repo.Get(&user, "SELECT * FROM users WHERE id = $1", userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetUser error")
	}
	return &user, nil
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}
