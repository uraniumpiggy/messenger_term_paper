package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/internal/users"
)

type db struct {
	*sql.DB
}

func NewStorage(database *sql.DB) users.Storage {
	return &db{database}
}

func (d *db) RegisterUser(ctx context.Context, data *users.UserRegisterRequest) error {
	var id int
	err := d.QueryRowContext(ctx, `select id from users where login = $1`, data.Login).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if id != 0 {
		return fmt.Errorf("User already exists")
	}
	res, err := d.ExecContext(ctx, `insert into users (login, username, password) values ($1, $2, $3)`, data.Login, data.Username, data.Password)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if affected == 0 || err != nil {
		return fmt.Errorf("err in db")
	}
	return nil

}

func (d *db) AuthUser(ctx context.Context, data *users.UserLoginRequest) (string, error) {
	var password string
	err := d.QueryRowContext(ctx, `select password from users where login = $1`, data.Login).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("User not found")
		}
		return "", err
	}
	return password, nil
}
