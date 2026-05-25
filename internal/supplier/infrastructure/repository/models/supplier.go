package models

import (
	"time"

	"gorm.io/gorm"
)

type Supplier struct {
	ID              uint           `gorm:"primaryKey"`
	Name            string         `gorm:"not null"`
	Status          string         `gorm:"not null;default:active"`
	CreatedByUserID *uint          `gorm:"index"`
	UpdatedByUserID *uint          `gorm:"index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
