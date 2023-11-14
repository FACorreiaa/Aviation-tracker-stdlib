package user

import (
	"github.com/FACorreiaa/go-ollama/internal/api/structs"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{db: db, redis: redis}
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
