package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					Name:     "testdb",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port: "8080",
				},
			},
			wantErr: false,
		},
		{
			name: "missing DB_HOST",
			cfg: Config{
				DB: DBConfig{
					Host:     "",
					Port:     5432,
					User:     "user",
					Password: "pass",
					Name:     "testdb",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port: "8080",
				},
			},
			wantErr: true,
			errMsg:  "DB_HOST is required",
		},
		{
			name: "missing DB_PASSWORD",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "",
					Name:     "testdb",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port: "8080",
				},
			},
			wantErr: true,
			errMsg:  "DB_PASSWORD is required",
		},
		{
			name: "missing DB_SSLMODE",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					Name:     "testdb",
					SSLMode:  "",
				},
				Server: ServerConfig{
					Port: "8080",
				},
			},
			wantErr: true,
			errMsg:  "DB_SSLMODE is required",
		},
		{
			name: "missing DB_USER",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "",
					Password: "pass",
					Name:     "testdb",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port: "8080",
				},
			},
			wantErr: true,
			errMsg:  "DB_USER is required",
		},
		{
			name: "missing DB_NAME",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					Name:     "",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port: "8080",
				},
			},
			wantErr: true,
			errMsg:  "DB_NAME is required",
		},
		{
			name: "missing SERVER_PORT",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					Name:     "testdb",
					SSLMode:  "disable",
				},
				Server: ServerConfig{
					Port: "",
				},
			},
			wantErr: true,
			errMsg:  "SERVER_PORT is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		expected string
	}{
		{
			name: "standard DSN",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "pass",
					Name:     "testdb",
					SSLMode:  "disable",
				},
			},
			expected: "host=localhost port=5432 dbname=testdb user=user password=pass sslmode=disable",
		},
		{
			name: "custom DSN",
			cfg: Config{
				DB: DBConfig{
					Host:     "192.168.1.1",
					Port:     5433,
					User:     "admin",
					Password: "secret",
					Name:     "prod",
					SSLMode:  "require",
				},
			},
			expected: "host=192.168.1.1 port=5433 dbname=prod user=admin password=secret sslmode=require",
		},
		{
			name: "empty password allowed in DSN",
			cfg: Config{
				DB: DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "",
					Name:     "testdb",
					SSLMode:  "disable",
				},
			},
			expected: "host=localhost port=5432 dbname=testdb user=user password= sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cfg.GetDSN())
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		setup       func()
		teardown    func()
		wantErr     bool
		errContains string
		validate    func(*Config, *testing.T)
	}{
		{
			name: "success load from env",
			setup: func() {
				viper.Set("DB_HOST", "localhost")
				viper.Set("DB_PORT", 5432)
				viper.Set("DB_USER", "user")
				viper.Set("DB_PASSWORD", "pass")
				viper.Set("DB_NAME", "testdb")
				viper.Set("DB_SSLMODE", "disable")
				viper.Set("SERVER_HOST", "0.0.0.0")
				viper.Set("SERVER_PORT", "8080")
				viper.Set("REDIS_HOST", "redis")
				viper.Set("REDIS_PORT", "6379")
				viper.Set("REDIS_PASSWORD", "redispass")
				viper.Set("LOG_MODE", "debug")
			},
			teardown: func() {
				viper.Reset()
			},
			wantErr: false,
			validate: func(cfg *Config, t *testing.T) {
				assert.Equal(t, "localhost", cfg.DB.Host)
				assert.Equal(t, 5432, cfg.DB.Port)
				assert.Equal(t, "user", cfg.DB.User)
				assert.Equal(t, "pass", cfg.DB.Password)
				assert.Equal(t, "testdb", cfg.DB.Name)
				assert.Equal(t, "disable", cfg.DB.SSLMode)
				assert.Equal(t, "0.0.0.0", cfg.Server.Host)
				assert.Equal(t, "8080", cfg.Server.Port)
				assert.Equal(t, "redis", cfg.Redis.Host)
				assert.Equal(t, "6379", cfg.Redis.Port)
				assert.Equal(t, "redispass", cfg.Redis.Password)
				assert.Equal(t, "debug", cfg.LogMode)

				expectedDSN := "host=localhost port=5432 dbname=testdb user=user password=pass sslmode=disable"
				assert.Equal(t, expectedDSN, cfg.GetDSN())
			},
		},
		{
			name: "missing required fields",
			setup: func() {
				viper.Set("DB_HOST", "localhost")
				viper.Set("DB_PORT", 5432)
				viper.Set("DB_USER", "user")
			},
			teardown:    func() { viper.Reset() },
			wantErr:     true,
			errContains: "DB_PASSWORD is required",
		},
		{
			name: "file not found but env works",
			setup: func() {
				viper.SetConfigName("nonexistent_config")
				viper.SetConfigType("env")
				viper.AddConfigPath(".")

				viper.Set("DB_HOST", "localhost")
				viper.Set("DB_PORT", 5432)
				viper.Set("DB_USER", "user")
				viper.Set("DB_PASSWORD", "pass")
				viper.Set("DB_NAME", "testdb")
				viper.Set("DB_SSLMODE", "disable")
				viper.Set("SERVER_PORT", "8080")
			},
			teardown: func() {
				viper.Reset()
			},
			wantErr: false,
			validate: func(cfg *Config, t *testing.T) {
				assert.Equal(t, "localhost", cfg.DB.Host)
				assert.Equal(t, "8080", cfg.Server.Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.teardown != nil {
				defer tt.teardown()
			}

			cfg, err := LoadConfig()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if tt.validate != nil {
					tt.validate(cfg, t)
				}
			}
		})
	}
}
