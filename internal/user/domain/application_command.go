package domain

import (
	"context"
	"errors"

	"golang-fiber/internal/user/domain/entity"
)

var (
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type RegisterInput struct {
	Username  string
	Email     string
	Password  string
	Firstname string
	Lastname  string
	Phone     string
}

type LoginInput struct {
	Username string
	Password string
}

type LoginResult struct {
	Token string
	User  *entity.User
}

type UserApplicationCommand interface {
	Register(context.Context, RegisterInput) (*entity.User, error)
	Login(context.Context, LoginInput) (*LoginResult, error)
}
