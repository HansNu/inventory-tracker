package service

import (
	"context"
	"errors"

	"inventory-tracker/apperr"
	"inventory-tracker/models"
	"inventory-tracker/repository"
	"inventory-tracker/utils"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterReq) (models.Users, string, error)
	Login(ctx context.Context, req models.LoginReq) (models.Users, string, error)
	CurrentUser(ctx context.Context, userID int) (models.Users, error)
}

type authService struct {
	Repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{Repo: repo}
}

// Register is where the auth rules live: uniqueness, hashing, token issuance.
// None of it needs an HTTP request or a live database to be exercised, which is
// the whole reason this layer exists.
func (s *authService) Register(ctx context.Context, req models.RegisterReq) (models.Users, string, error) {
	exists, err := s.Repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return models.Users{}, "", err
	}
	if exists {
		return models.Users{}, "", apperr.ErrDuplicate
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.Users{}, "", err
	}

	user, err := s.Repo.Create(ctx, req.Name, req.CategoryGroup, req.Username, hash)
	if err != nil {
		return models.Users{}, "", err
	}

	token, err := utils.GenerateToken(user.Id, user.CategoryGroup)
	if err != nil {
		return models.Users{}, "", err
	}
	return user, token, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginReq) (models.Users, string, error) {
	user, err := s.Repo.GetByUsername(ctx, req.Username)
	if err != nil {
		// "No such user" and "wrong password" must be indistinguishable to the
		// caller, or the endpoint becomes a username oracle.
		if errors.Is(err, apperr.ErrNotFound) {
			return models.Users{}, "", apperr.ErrInvalidCredentials
		}
		return models.Users{}, "", err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return models.Users{}, "", apperr.ErrInvalidCredentials
	}

	// Never let the hash past this layer.
	user.Password = ""

	token, err := utils.GenerateToken(user.Id, user.CategoryGroup)
	if err != nil {
		return models.Users{}, "", err
	}
	return user, token, nil
}

func (s *authService) CurrentUser(ctx context.Context, userID int) (models.Users, error) {
	return s.Repo.GetByID(ctx, userID)
}
