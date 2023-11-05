package service

import (
	"github.com/FACorreiaa/go-ollama/internal/api/repository"
	"github.com/FACorreiaa/go-ollama/internal/api/service/user"
	"github.com/FACorreiaa/go-ollama/internal/api/structs"
)

type User interface {
	Create(user structs.User) (id int, err error)
	GetById(id int) (user structs.User, err error)
	GetByUsername(username string) (user structs.User, err error)
}

type Service struct {
	User User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: user.NewService(repo),
	}
}
