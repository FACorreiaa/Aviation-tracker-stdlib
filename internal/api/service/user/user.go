package user

import (
	"github.com/FACorreiaa/go-ollama/internal/api/repository"
	"github.com/FACorreiaa/go-ollama/internal/api/structs"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(user structs.User) (int, error) {
	return s.repo.User.Create(user)
}

func (s *Service) GetById(id int) (structs.User, error) {
	return s.repo.User.GetById(id)
}

func (s *Service) GetByUsername(username string) (structs.User, error) {
	return s.repo.User.GetByUsername(username)
}
