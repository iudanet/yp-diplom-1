package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		cfg := New()
		if cfg.HttpAddress != "localhost:8080" {
			t.Errorf("expected HttpAddress 'localhost:8080', got '%s'", cfg.HttpAddress)
		}
		if cfg.DatabaseUri != "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" {
			t.Errorf("unexpected default DatabaseUri: %s", cfg.DatabaseUri)
		}
		if cfg.AccrualSystemAddress != "http://localhost:8081" {
			t.Errorf("unexpected default AccrualSystemAddress: %s", cfg.AccrualSystemAddress)
		}
	})

	t.Run("flag parsing", func(t *testing.T) {
		// Сохраняем оригинальные аргументы
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// Устанавливаем тестовые аргументы
		os.Args = []string{
			"cmd",
			"-a", "localhost:9090",
			"-d", "postgres://user:pass@localhost:5432/testdb",
			"-r", "http://accrual:8082",
		}

		cfg := New()
		cfg.FlagParse()

		if cfg.HttpAddress != "localhost:9090" {
			t.Errorf("expected HttpAddress 'localhost:9090', got '%s'", cfg.HttpAddress)
		}
		if cfg.DatabaseUri != "postgres://user:pass@localhost:5432/testdb" {
			t.Errorf("unexpected DatabaseUri: %s", cfg.DatabaseUri)
		}
		if cfg.AccrualSystemAddress != "http://accrual:8082" {
			t.Errorf("unexpected AccrualSystemAddress: %s", cfg.AccrualSystemAddress)
		}
	})

	t.Run("env parsing", func(t *testing.T) {
		// Сохраняем оригинальные переменные окружения
		oldHttpAddr := os.Getenv("RUN_ADDRESS")
		oldDbUri := os.Getenv("DATABASE_URI")
		oldAccrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
		defer func() {
			os.Setenv("RUN_ADDRESS", oldHttpAddr)
			os.Setenv("DATABASE_URI", oldDbUri)
			os.Setenv("ACCRUAL_SYSTEM_ADDRESS", oldAccrualAddr)
		}()

		// Устанавливаем тестовые переменные окружения
		os.Setenv("RUN_ADDRESS", "localhost:9091")
		os.Setenv("DATABASE_URI", "postgres://envuser:envpass@localhost:5432/envdb")
		os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://envaccrual:8083")

		cfg := New()
		cfg.EnvParse()

		if cfg.HttpAddress != "localhost:9091" {
			t.Errorf("expected HttpAddress 'localhost:9091', got '%s'", cfg.HttpAddress)
		}
		if cfg.DatabaseUri != "postgres://envuser:envpass@localhost:5432/envdb" {
			t.Errorf("unexpected DatabaseUri: %s", cfg.DatabaseUri)
		}
		if cfg.AccrualSystemAddress != "http://envaccrual:8083" {
			t.Errorf("unexpected AccrualSystemAddress: %s", cfg.AccrualSystemAddress)
		}
	})

	t.Run("string representation", func(t *testing.T) {
		cfg := New()
		expected := "HttpAddress: localhost:8080, DatabaseUri: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable, AccrualSystemAddress: http://localhost:8081"
		if cfg.String() != expected {
			t.Errorf("expected string '%s', got '%s'", expected, cfg.String())
		}
	})
}
