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
}

func New() *Config {
	return &Config{
		HTTPAddress:          "localhost:8080",
		DatabaseURI:          "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		AccrualSystemAddress: "http://localhost:8081",
	}
}

func (c *Config) FlagParse() {

	flag.StringVar(&c.HTTPAddress, "a", c.HTTPAddress, "http address| env RUN_ADDRESS ")
	flag.StringVar(&c.DatabaseURI, "d", c.DatabaseURI, "database uri | env DATABASE_URI ")
	flag.StringVar(&c.AccrualSystemAddress, "r", c.AccrualSystemAddress, "accrual system address | env ACCRUAL_SYSTEM_ADDRESS ")

	flag.Parse()
}
func (c *Config) EnvParse() {
	c.HTTPAddress = os.Getenv("RUN_ADDRESS")
	c.DatabaseURI = os.Getenv("DATABASE_URI")
	c.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
}

func (c *Config) String() string {
	return fmt.Sprintf("HTTPAddress: %s, DatabaseURI: %s, AccrualSystemAddress: %s", c.HTTPAddress, c.DatabaseURI, c.AccrualSystemAddress)
}
