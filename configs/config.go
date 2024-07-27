package configs

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time_tracker/internal/repository"
)

func LoadMain() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
		panic(err)
	}
	logrus.Info("Configuration file loaded")
}

func GetDBConfig() repository.DbConfig {
	return repository.DbConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Name:     viper.GetString("DB_NAME"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
	}
}
