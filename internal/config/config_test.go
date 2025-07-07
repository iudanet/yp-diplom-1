package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		teardown func()
		config   *Config
		check    func(t *testing.T, cfg *Config)
	}{
		{
			name:   "default_values",
			config: New(),
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "localhost:8080", cfg.HTTPAddress)
				assert.Equal(t, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", cfg.DatabaseURI)
				assert.Equal(t, "http://localhost:8081", cfg.AccrualSystemAddress)
			},
		},
		{
			name: "flag_parsing",
			setup: func() {
				oldArgs := os.Args
				os.Args = []string{
					"cmd",
					"-a", "localhost:1010",
					"-d", "postgres://user:pass@localhost66:5432/testdb",
					"-r", "http://accrual:8082",
				}
				t.Cleanup(func() { os.Args = oldArgs })
			},
			config: New(),
			check: func(t *testing.T, cfg *Config) {
				cfg.FlagParse()
				assert.Equal(t, "localhost:1010", cfg.HTTPAddress)
				assert.Equal(t, "postgres://user:pass@localhost66:5432/testdb", cfg.DatabaseURI)
				assert.Equal(t, "http://accrual:8082", cfg.AccrualSystemAddress)
			},
		},
		{
			name: "env_parsing",
			setup: func() {
				os.Setenv("RUN_ADDRESS", "localhost:9091")
				os.Setenv("DATABASE_URI", "postgres://envuser:envpass@localhost:5432/envdb")
				os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://envaccrual:8083")
			},
			teardown: func() {
				os.Unsetenv("RUN_ADDRESS")
				os.Unsetenv("DATABASE_URI")
				os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
			},
			config: New(),
			check: func(t *testing.T, cfg *Config) {
				cfg.EnvParse()
				assert.Equal(t, "localhost:9091", cfg.HTTPAddress)
				assert.Equal(t, "postgres://envuser:envpass@localhost:5432/envdb", cfg.DatabaseURI)
				assert.Equal(t, "http://envaccrual:8083", cfg.AccrualSystemAddress)
			},
		},
		{
			name:   "string_representation",
			config: New(),
			check: func(t *testing.T, cfg *Config) {
				expected := "HTTPAddress: localhost:8080, DatabaseURI: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable, AccrualSystemAddress: http://localhost:8081"
				assert.Equal(t, expected, cfg.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.teardown != nil {
				t.Cleanup(tt.teardown)
			}

			tt.check(t, tt.config)
		})
	}
}
