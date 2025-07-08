package config

import (
	"flag"
	"fmt"
	"os"
)

// - адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
// - адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
// - адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
type Config struct {
	HTTPAddress string
	// dsn "postgres://postgres:postgres@localhost:5432/metrics_db?sslmode=disable"
	DatabaseURI          string
	AccrualSystemAddress string
	Migrate              bool
}

func New() *Config {
	return &Config{
		HTTPAddress:          "localhost:8080",
		DatabaseURI:          "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		AccrualSystemAddress: "http://localhost:8081",
		Migrate:              true,
	}
}

func (c *Config) FlagParse() {
	flag.StringVar(&c.HTTPAddress, "a", c.HTTPAddress, "http address| env RUN_ADDRESS ")
	flag.StringVar(&c.DatabaseURI, "d", c.DatabaseURI, "database uri | env DATABASE_URI ")
	flag.StringVar(
		&c.AccrualSystemAddress,
		"r",
		c.AccrualSystemAddress,
		"accrual system address | env ACCRUAL_SYSTEM_ADDRESS ",
	)
	flag.BoolVar(
		&c.Migrate,
		"m",
		c.Migrate,
		"migrate database | env MIGRATE ",
	)

	flag.Parse()
}

func (c *Config) EnvParse() {
	if os.Getenv("RUN_ADDRESS") != "" {
		c.HTTPAddress = os.Getenv("RUN_ADDRESS")
	}
	if os.Getenv("DATABASE_URI") != "" {
		c.DatabaseURI = os.Getenv("DATABASE_URI")
	}
	if os.Getenv("ACCRUAL_SYSTEM_ADDRESS") != "" {
		c.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	}
	if os.Getenv("MIGRATE") != "" {
		c.Migrate = os.Getenv("MIGRATE") == "true"
	}
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"Config:\n HTTPAddress: %s \n DatabaseURI: %s\n AccrualSystemAddress: %s\n",
		c.HTTPAddress,
		c.DatabaseURI,
		c.AccrualSystemAddress,
	)
}
