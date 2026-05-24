package models

import "time"

type Supplier struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Supplier) TableName() string {
	return "suppliers"
}
