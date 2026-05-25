package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host     string `split_words:"true" default:"localhost"`
	Port     string `split_words:"true" default:"5432"`
	User     string `split_words:"true" default:"postgres"`
	Password string `split_words:"true" default:"postgres"`
	DBName   string `split_words:"true" default:"product_db"`
	SSLMode  string `split_words:"true" default:"disable"`
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Bangkok",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

func (d DBConfig) Connect(debug bool) (*gorm.DB, error) {
	logLevel := logger.Silent
	if debug {
		logLevel = logger.Info
	}
	return gorm.Open(postgres.Open(d.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
}
