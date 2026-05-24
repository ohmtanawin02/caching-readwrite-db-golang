package entity

import "time"

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	Firstname string
	Lastname  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
