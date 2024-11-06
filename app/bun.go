package app

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type BunDriver struct {
	db *bun.DB
}

func newBunDriver(dsn string) (*BunDriver, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &BunDriver{db: db}, nil
}

func (d *BunDriver) InsertUser(ctx context.Context, user User) error {
	_, err := d.db.NewInsert().Model(&user).
		Column("username", "password", "city").
		Exec(ctx)
	return err
}

func (d *BunDriver) SelectUser(ctx context.Context, username string) (User, error) {
	var user User
	err := d.db.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx, &user)
	return user, err
}

func (d *BunDriver) Disconnect() error {
	return d.db.Close()
}
