package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Username  string         `gorm:"not null"`
	Email     string         `gorm:"not null"`
	Password  string         `gorm:"not null"`
	Firstname string         `gorm:"not null"`
	Lastname  string         `gorm:"not null"`
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
