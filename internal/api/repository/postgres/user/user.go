package user

import (
	"github.com/FACorreiaa/go-ollama/internal/api/structs"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user structs.User) (int, error) {
	return 1, nil
}

func (r *Repository) GetById(id int) (structs.User, error) {
	return structs.User{}, nil
}

func (r *Repository) GetByUsername(username string) (structs.User, error) {
	return structs.User{}, nil
}
