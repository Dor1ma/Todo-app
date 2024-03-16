package service

import (
	"Todo-app/internal/models"
	"Todo-app/internal/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user *models.User) (*models.User, error) {
	return s.repo.Create(user)
}

func (s *AuthService) FindByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *AuthService) Find(id int) (*models.User, error) {
	return s.repo.Find(id)
}
