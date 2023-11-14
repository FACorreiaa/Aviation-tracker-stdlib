package account

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

func (a *Accounts) Login(ctx context.Context, form LoginForm) (*Token, error) {
	if err := a.validator.Struct(form); err != nil {
		return nil, err
	}

	rows, _ := a.pgpool.Query(
		ctx,
		`
		select
			user_id,
			username,
			email,
			password_hash,
			bio,
			image,
			created_at,
			updated_at
		from "user" where email = $1 limit 1
		`,
		form.Email,
	)
	user, err := pgx.CollectOneRow[User](rows, pgx.RowToStructByPos[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}

		slog.Error("Error querying user", "err", err)
		return nil, errors.New("internal server error")
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(form.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	tokenBytes := make([]byte, RAND_SIZE)
	if _, err := rand.Read(tokenBytes); err != nil {
		slog.Error("Error generating token", "err", err)
		return nil, errors.New("internal server error")
	}

	token := Token(fmt.Sprintf("%x", tokenBytes))

	if _, err := a.pgpool.Exec(
		ctx,
		`
		insert into user_token (user_id, token, context)
		values ($1, $2, $3)
		`,
		user.ID,
		token,
		"auth",
	); err != nil {
		slog.Error("Error inserting token", "err", err)
		return nil, errors.New("internal server error")
	}

	return &token, nil
}

func (h *Accounts) UserFromSessionToken(ctx context.Context, token Token) (*User, error) {
	rows, _ := h.pgpool.Query(
		ctx,
		`
		select
			t.token,
			t.created_at,
			row(
				u.user_id,
				u.username,
				u.email,
				u.password_hash,
				u.bio,
				u.image,
				u.created_at,
				u.updated_at
			)
		from "user" u
		join user_token t using (user_id)
		where token = $1 and context = 'auth'
		limit 1
		`,
		token,
	)
	userWithToken, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[UserToken])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("auth session expired")
		}

		slog.Error("Error querying user with token", "err", err)
		return nil, errors.New("internal server error")
	}

	if userWithToken.CreatedAt == nil || time.Since(*userWithToken.CreatedAt) > MAX_AGE {
		return nil, errors.New("auth session expired")
	}

	return userWithToken.User, nil
}
