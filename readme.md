# PostgreSQL Drivers Benchmark

This is a benchmarking tool for comparing different PostgreSQL drivers in Go. The application measures and compares the performance of popular PostgreSQL drivers: pq, pgx, GORM, and Bun.

## PostgreSQL Drivers Benchmark

- lib/pq: Pure Go PostgreSQL driver
- jackc/pgx: High-performance PostgreSQL driver
- gorm: Popular ORM with PostgreSQL support
- uptrace/bun: Lightweight PostgreSQL ORM

## Build

docker-compose build: (only postgres)

```
docker-compose up --build
```

DB migration: https://github.com/golang-migrate/migrate

.env:

```
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_USERNAME=
DB_DBNAME=
DB_SSLMODE=disable
```

config.yaml structure:

```
pg_bench:
  goroutines_pool: 95     # Number of concurrent workers
  requests_amount: 10000   # Total number of requests to perform
  timeout: 5000ms           # Benchmark timeout duration
```

Example output:

```
Driver    | Action    | RPS       | Total requests | Err amount
----------|-----------|-----------|----------------|------------
pq        | Insert    | 5234.21   | 10000         | 0
pq        | Select    | 6123.45   | 10000         | 0
----------|-----------|-----------|----------------|------------
pgx       | Insert    | 7234.21   | 10000         | 0
pgx       | Select    | 8123.45   | 10000         | 0
----------|-----------|-----------|----------------|------------
bun       | Insert    | 6234.21   | 10000         | 0
bun       | Select    | 7123.45   | 10000         | 0
----------|-----------|-----------|----------------|------------
gorm      | Insert    | 4234.21   | 10000         | 0
gorm      | Select    | 5123.45   | 10000         | 0
```
