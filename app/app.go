package app

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kolibriee/bench-db-comparison/app/config"
)

type User struct {
	ID           int       `bun:"id,pk,autoincrement" gorm:"column:id"`
	Username     string    `bun:"username" gorm:"column:username"`
	Password     string    `bun:"password" gorm:"column:password"`
	City         string    `bun:"city" gorm:"column:city"`
	RegisteredAt time.Time `bun:"registered_at" gorm:"column:registered_at"`
}

type DBDriver interface {
	InsertUser(ctx context.Context, user User) error
	SelectUser(ctx context.Context, username string) (User, error)
	Disconnect() error
}

type BenchmarkResult struct {
	Rps           float64
	TotalRequests int
	ErrAmount     int
}

func Run(configPath string, fileName string) {
	cfg, err := config.New(configPath, fileName)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	testData := genData(cfg.BenchmarkConfig.RequestsAmount)
	pqDriver, err := newPQDriver(dsn)
	if err != nil {
		log.Fatalf("failed to create pq driver: %v", err)
	}
	pgxDriver, err := newPgxDriver(dsn)
	if err != nil {
		log.Fatalf("failed to create pgx driver: %v", err)
	}

	bunDriver, err := newBunDriver(dsn)
	if err != nil {
		log.Fatalf("failed to create bun driver: %v", err)
	}

	gormDriver, err := newGormDriver(dsn)
	if err != nil {
		log.Fatalf("failed to create gorm driver: %v", err)
	}
	drivers := []struct {
		name   string
		driver DBDriver
	}{
		{"pq", pqDriver},
		{"pgx", pgxDriver},
		{"bun", bunDriver},
		{"gorm", gormDriver},
	}
	if err := MigrateDown(dsn); err != nil {
		log.Fatalf("failed to migrate down: %v", err)
	}
	fmt.Printf("%-10s | %-10s | %-10s | %-15s | %-10s\n", "Driver", "Action", "RPS", "Total requests", "Err amount")
	for _, driver := range drivers {
		time.Sleep(time.Second * 3)
		if err := MigrateUp(dsn); err != nil {
			log.Fatalf("failed to migrate up: %v", err)
		}

		insertResult := runBenchmark("insert", driver.driver, testData, cfg.BenchmarkConfig)
		fmt.Printf("%-10s | %-10s | %-10.2f | %-15d | %-10d\n",
			driver.name, "Insert", insertResult.Rps, insertResult.TotalRequests, insertResult.ErrAmount)

		selectResult := runBenchmark("select", driver.driver, testData, cfg.BenchmarkConfig)
		fmt.Printf("%-10s | %-10s | %-10.2f | %-15d | %-10d\n",
			driver.name, "Select", selectResult.Rps, selectResult.TotalRequests, selectResult.ErrAmount)

		if err := driver.driver.Disconnect(); err != nil {
			log.Fatalf("failed to disconnect: %v", err)
		}

		if err := MigrateDown(dsn); err != nil {
			log.Fatalf("failed to migrate down: %v", err)
		}
		fmt.Println(strings.Repeat("-", 70))
	}
}
