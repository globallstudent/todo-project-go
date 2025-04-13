package services

import (
	"context"
	"errors"

	"github.com/globallstudent/todo-project-go/internal/models"
	"github.com/globallstudent/todo-project-go/internal/repositories"
	"github.com/globallstudent/todo-project-go/internal/utils"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, username, password, role string) (*models.User, error) {

	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	_, err := s.userRepo.FindUserByUsername(ctx, username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.FindUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
