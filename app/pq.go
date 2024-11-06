package app

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type PQDriver struct {
	db *sql.DB
}

func newPQDriver(dsn string) (*PQDriver, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PQDriver{db: db}, nil
}

func (d *PQDriver) InsertUser(ctx context.Context, user User) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO users (username, password, city) VALUES ($1, $2, $3)", user.Username, user.Password, user.City)
	return err
}

func (d *PQDriver) SelectUser(ctx context.Context, username string) (User, error) {
	var user User
	err := d.db.QueryRowContext(ctx, "SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.City, &user.RegisteredAt)
	return user, err
}

func (d *PQDriver) Disconnect() error {
	return d.db.Close()
}
