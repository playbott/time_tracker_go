package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time_tracker/internal/models"
)

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func NewPostgresDB(dbConfig DbConfig) (*gorm.DB, error) {
	err := createDatabase(dbConfig)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	logrus.Infof("Database connected. DSN: %v", fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Name,
		dbConfig.Port,
	))

	autoMigrate := viper.GetBool("DB_AUTO_MIGR")
	if autoMigrate {
		logrus.Info("Database automigration...")
		err = db.AutoMigrate(&models.User{}, &models.Task{})
		if err != nil {
			panic(err)
		}
		logrus.Info("Database migration completed")
	}

	return db, nil
}

func createDatabase(dbConfig DbConfig) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	exists := false
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbConfig.Name)
	err = db.Raw(query).Scan(&exists).Error
	if err != nil {
		logrus.Errorf("failed to check if database exists: %w", err)
		return err
	}

	if !exists {
		query = fmt.Sprintf("CREATE DATABASE %s", dbConfig.Name)
		err = db.Exec(query).Error
		if err != nil {
			logrus.Errorf("failed to create database: %w", err)
			return err
		}
		logrus.Printf("Database %s created successfully", dbConfig.Name)
	}

	return nil
}
