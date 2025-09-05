package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigFile("./config.env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
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

func GetServerPort() string {
	return viper.GetString("SERVER_PORT")
}

func GetServerHost() string {
	return viper.GetString("SERVER_HOST")
}

func GetEnv() string {
	return viper.GetString("ENV")
}
