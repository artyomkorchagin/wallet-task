package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB       DBConfig
	Server   ServerConfig
	Redis    RedisConfig
	LogLevel string `mapstructure:"LOG_LEVEL"`
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
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	viper.AutomaticEnv()

	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("LOG_LEVEL", "info")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func GetDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_SSLMODE"))
}

func (cfg *Config) Validate() error {
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
