package app

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDriver struct {
	db *gorm.DB
}

func newGormDriver(dsn string) (*GormDriver, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return &GormDriver{db: db}, nil
}

func (d *GormDriver) InsertUser(ctx context.Context, user User) error {
	return d.db.WithContext(ctx).Select("username", "password", "city").Create(&user).Error
}

func (d *GormDriver) SelectUser(ctx context.Context, username string) (User, error) {
	var user User
	err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return user, err
}

func (d *GormDriver) Disconnect() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
