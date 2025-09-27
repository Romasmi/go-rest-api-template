package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Romasmi/go-rest-api-template/internal/utils"
	"strconv"

	"github.com/Romasmi/go-rest-api-template/internal/middleware"
	"github.com/Romasmi/go-rest-api-template/internal/models"
	"github.com/Romasmi/go-rest-api-template/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, user *models.UserCreate) (*models.User, error) {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, id int, user *models.UserUpdate) (*models.User, error) {
	return s.repo.Update(ctx, id, user)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context, page, pageSize int) ([]*models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *UserService) Login(ctx context.Context, login *models.UserLogin) (string, error) {
	user, err := s.repo.GetByUsername(ctx, login.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", fmt.Errorf("invalid username or password")
		}
		return "", err
	}

	if !utils.CheckPassword(login.Password, user.PasswordHash) {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := middleware.GenerateJWT(strconv.Itoa(user.ID), user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *UserService) Register(ctx context.Context, user *models.UserCreate) (string, error) {
	newUser, err := s.repo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return "", fmt.Errorf("username or email already exists")
		}
		return "", err
	}

	token, err := middleware.GenerateJWT(strconv.Itoa(newUser.ID), newUser.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
