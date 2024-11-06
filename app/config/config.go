package config

import (
	"errors"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	BenchmarkConfig `mapstructure:"pg_bench"`
	Postgres
}
type BenchmarkConfig struct {
	GoroutinesPool int           `mapstructure:"goroutines_pool"`
	RequestsAmount int           `mapstructure:"requests_amount"`
	Timeout        time.Duration `mapstructure:"timeout"`
}

type Postgres struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func New(path string, fileName string) (*Config, error) {
	var cfg *Config

	godotenv.Load()

	viper.SetConfigName(fileName)
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New("failed to read config file: " + err.Error())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.New("failed to unmarshal config: " + err.Error())
	}

	if err := envconfig.Process("DB", &cfg.Postgres); err != nil {
		return nil, errors.New("failed to process env variables: " + err.Error())
	}

	return cfg, nil
}
