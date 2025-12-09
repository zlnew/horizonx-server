package domain

import (
	"context"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
}

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Register(ctx context.Context, req RegisterRequest) error
}
