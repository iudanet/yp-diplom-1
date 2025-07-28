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
			name:   "defaultValues",
			config: New(),
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "localhost:8080", cfg.HTTPAddress)
				assert.Equal(
					t,
					"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
					cfg.DatabaseURI,
				)
				assert.Equal(t, "http://localhost:8081", cfg.AccrualSystemAddress)
			},
		},
		{
			name: "flagParsing",
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
			name: "envParsing",
			setup: func() {
				err := os.Setenv("RUN_ADDRESS", "localhost:9091")
				assert.NoError(t, err)
				err = os.Setenv("DATABASE_URI", "postgres://envuser:envpass@localhost:5432/envdb")
				assert.NoError(t, err)
				err = os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://envaccrual:8083")
				assert.NoError(t, err)
			},
			teardown: func() {
				err := os.Unsetenv("RUN_ADDRESS")
				assert.NoError(t, err)
				err = os.Unsetenv("DATABASE_URI")
				assert.NoError(t, err)
				err = os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
				assert.NoError(t, err)
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
			name:   "stringRepresentation",
			config: New(),
			check: func(t *testing.T, cfg *Config) {
				// expected := "HTTPAddress: localhost:8080, DatabaseURI: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable, AccrualSystemAddress: http://localhost:8081"
				expected := "Config:\n HTTPAddress: localhost:8080 \n DatabaseURI: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable\n AccrualSystemAddress: http://localhost:8081\n"
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
