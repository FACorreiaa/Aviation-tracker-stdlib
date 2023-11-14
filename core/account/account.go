package account

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const RAND_SIZE = 32
const MAX_AGE = time.Hour * 24 * 60

type Token = string

type Accounts struct {
	pgpool    *pgxpool.Pool
	validator *validator.Validate
}

func NewAccounts(
	pgpool *pgxpool.Pool,
	validator *validator.Validate,
) *Accounts {
	return &Accounts{
		pgpool:    pgpool,
		validator: validator,
	}
}

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash []byte
	Bio          string
	Image        *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type UserToken struct {
	Token     string
	CreatedAt *time.Time
	User      *User
}

func (h *Accounts) Logout(ctx context.Context, token Token) error {
	if _, err := h.pgpool.Exec(
		ctx,
		`
		delete from user_token where token_id = $1 and context = 'auth'
		`,
		token,
	); err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}

	return nil
}
