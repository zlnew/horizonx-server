// Package auth
package auth

import (
	"context"

	"horizonx-server/internal/domain"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User        *domain.User `json:"user"`
	AccessToken string       `json:"access_token"`
}

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Register(ctx context.Context, req RegisterRequest) error
}
