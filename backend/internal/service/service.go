package service

import (
	"context"

	"auth/service/internal/logger"
	"auth/service/internal/model"
)

var (
	log = logger.L()
)

type UserRepoService interface {
	Create(ctx context.Context, user model.Model) (int64, error)
	GetByEmail(ctx context.Context, email string) (*model.Model, error)
	GetByID(ctx context.Context, id int64) (*model.Model, error)
}

type UserService struct {
	repo UserRepoService
}

func NewUserService(repo UserRepoService) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) Create(ctx context.Context, user model.Model) (int64, error) {
	id, err := u.repo.Create(ctx, user)
	if err != nil {
		log.Errorf("Failed: %v", err)
		return 0, err
	}

	return id, nil
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (*model.Model, error) {
	task, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Errorf("Failed: %v", err)
		return nil, err
	}

	return task, nil
}

func (u *UserService) GetByID(ctx context.Context, id int64) (*model.Model, error) {
	task, err := u.repo.GetByID(ctx, id)
	if err != nil {
		log.Errorf("Failed: %v", err)
		return nil, err
	}
	return task, nil
}
