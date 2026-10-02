package services

import (
	"errors"

	"github.com/raifuki/task-management/internal/models"
	"github.com/raifuki/task-management/internal/repositories"
	"github.com/raifuki/task-management/pkg/hash"
	"github.com/raifuki/task-management/pkg/jwt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	jwtMgr   *jwt.JWTManager
}

func NewAuthService(u *repositories.UserRepository, j *jwt.JWTManager) *AuthService {
	return &AuthService{userRepo: u, jwtMgr: j}
}

type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *AuthService) Register(in RegisterInput) (*models.User, error) {
	existing, err := s.userRepo.FindByEmail(in.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := hash.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: hashed,
		Role:     models.RoleUser,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(in LoginInput) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(in.Email)
	if err != nil {
		return "", nil, err
	}
	if user == nil || !hash.CheckPassword(user.Password, in.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := s.jwtMgr.Generate(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}