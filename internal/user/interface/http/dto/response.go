package dto

import (
	"golang-fiber/internal/user/domain"
	"golang-fiber/internal/user/domain/entity"
)

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func ToUserResponse(u *entity.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Phone:     u.Phone,
	}
}

func ToLoginResponse(result *domain.LoginResult) LoginResponse {
	return LoginResponse{
		Token: result.Token,
		User:  ToUserResponse(result.User),
	}
}
