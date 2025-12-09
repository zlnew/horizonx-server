// Package user
package user

import (
	"context"

	"horizonx-server/internal/config"
	"horizonx-server/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo domain.UserRepository
	cfg  *config.Config
}

func NewService(repo domain.UserRepository, cfg *config.Config) domain.UserService {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Get(ctx context.Context, opts domain.ListOptions) (*domain.ListResult[*domain.User], error) {
	if opts.IsPaginate {
		if opts.Page <= 0 {
			opts.Page = 1
		}
		if opts.Limit <= 0 {
			opts.Limit = 10
		}
	}

	users, total, err := s.repo.GetUsers(ctx, opts)
	if err != nil {
		return nil, err
	}

	res := &domain.ListResult[*domain.User]{
		Data: users,
		Meta: nil,
	}

	if opts.IsPaginate {
		res.Meta = domain.CalculateMeta(total, opts.Page, opts.Limit)
	}

	return res, nil
}

func (s *service) Create(ctx context.Context, req domain.UserSaveRequest) error {
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

func (s *service) Update(ctx context.Context, req domain.UserSaveRequest, userID int64) error {
	if user, _ := s.repo.GetUserByEmail(ctx, req.Email); user != nil {
		if user.ID != userID {
			return domain.ErrEmailAlreadyExists
		}
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPwd),
	}

	return s.repo.UpdateUser(ctx, user, userID)
}

func (s *service) Delete(ctx context.Context, userID int64) error {
	if _, err := s.repo.GetUserByID(ctx, userID); err != nil {
		return err
	}

	return s.repo.DeleteUser(ctx, userID)
}
