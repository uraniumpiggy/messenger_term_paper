package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewClient(ctx context.Context, host, port, username, password, databaseName string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, databaseName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to database due to error: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Cannot ping database due to error: %s", err)
	}

	db.SetConnMaxIdleTime(time.Minute)
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)

	return db, nil
}
