package auth

import (
	"context"
	"time"

	"horizonx-server/internal/config"
	"horizonx-server/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo domain.UserRepository
	cfg  *config.Config
}

func NewService(repo domain.UserRepository, cfg *config.Config) AuthService {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) error {
	if user, _ := s.repo.GetUserByEmail(ctx, req.Email); user != nil {
		return domain.ErrEmailAlreadyExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPwd),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(s.cfg.JWTExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: tokenString,
		User:        user,
	}, nil
}
