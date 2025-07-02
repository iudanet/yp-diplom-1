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
	HttpAddress string
	// dsn "postgres://postgres:postgres@localhost:5432/metrics_db?sslmode=disable"
	DatabaseUri          string
	AccrualSystemAddress string
}

func New() *Config {
	return &Config{
		HttpAddress:          "localhost:8080",
		DatabaseUri:          "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		AccrualSystemAddress: "http://localhost:8081",
	}
}

func (c *Config) FlagParse() {

	flag.StringVar(&c.HttpAddress, "a", c.HttpAddress, "http address| env RUN_ADDRESS ")
	flag.StringVar(&c.DatabaseUri, "d", c.DatabaseUri, "database uri | env DATABASE_URI ")
	flag.StringVar(&c.AccrualSystemAddress, "r", c.AccrualSystemAddress, "accrual system address | env ACCRUAL_SYSTEM_ADDRESS ")

	flag.Parse()
}
func (c *Config) EnvParse() {
	c.HttpAddress = os.Getenv("RUN_ADDRESS")
	c.DatabaseUri = os.Getenv("DATABASE_URI")
	c.AccrualSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
}

func (c *Config) String() string {
	return fmt.Sprintf("HttpAddress: %s, DatabaseUri: %s, AccrualSystemAddress: %s", c.HttpAddress, c.DatabaseUri, c.AccrualSystemAddress)
}
