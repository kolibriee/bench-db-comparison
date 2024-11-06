package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxDriver struct {
	db *pgxpool.Pool
}

func newPgxDriver(dsn string) (*PgxDriver, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &PgxDriver{db: db}, nil
}

func (d *PgxDriver) InsertUser(ctx context.Context, user User) error {
	_, err := d.db.Exec(ctx, "INSERT INTO users (username, password, city) VALUES ($1, $2, $3)", user.Username, user.Password, user.City)
	return err
}

func (d *PgxDriver) SelectUser(ctx context.Context, username string) (User, error) {
	var user User
	err := d.db.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.City, &user.RegisteredAt)
	return user, err
}

func (d *PgxDriver) Disconnect() error {
	d.db.Close()
	return nil
}
