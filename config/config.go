package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB      DBConfig     `mapstructure:",squash"`
	Server  ServerConfig `mapstructure:",squash"`
	Redis   RedisConfig  `mapstructure:",squash"`
	LogMode string       `mapstructure:"LOG_MODE"`
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type ServerConfig struct {
	Host string `mapstructure:"SERVER_HOST"`
	Port string `mapstructure:"SERVER_PORT"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.SSLMode)
}

func (cfg *Config) Validate() error {
	if cfg.DB.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if cfg.DB.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if cfg.DB.SSLMode == "" {
		return fmt.Errorf("DB_SSLMODE is required")
	}
	if cfg.DB.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if cfg.DB.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if cfg.DB.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if cfg.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}
	return nil
}
