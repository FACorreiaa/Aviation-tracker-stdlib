package repository

import (
	"github.com/FACorreiaa/go-ollama/internal/api/repository/postgres"
	"github.com/FACorreiaa/go-ollama/internal/api/repository/postgres/user"
	"github.com/FACorreiaa/go-ollama/internal/api/structs"
)

type Config struct {
	postgresConfig postgres.Config
}

func NewConfig(postgresConfig postgres.Config) Config {
	return Config{postgresConfig: postgresConfig}
}

type User interface {
	Create(user structs.User) (id int, err error)
	GetById(id int) (user structs.User, err error)
	GetByUsername(username string) (user structs.User, err error)
}

type Repository struct {
	User User
}

func NewRepository(config Config) *Repository {
	psql := postgres.NewPostgres(config.postgresConfig)
	return &Repository{
		User: user.NewRepository(psql.GetDB()),
	}
}
