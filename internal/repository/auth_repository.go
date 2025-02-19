package repository

import (
	handle_errors "avitoTT/internal/errors"
	"avitoTT/openapi/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{DB: db}
}

func (r *AuthRepositoryImpl) Authenticate(req models.AuthRequest) error {
	var password string

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := r.DB.QueryRow(ctx, "SELECT password FROM users WHERE username = $1", req.Username).Scan(&password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.As(err, &pgx.ErrNoRows) {
			_, err = r.DB.Exec(ctx, `
    INSERT INTO users (username, password)
    VALUES ($1, $2)
    ON CONFLICT (username) DO NOTHING
`, req.Username, req.Password)

			if err != nil {
				return handle_errors.ErrDatabaseIssue
			}
			return nil
		}
		return handle_errors.ErrDatabaseIssue
	}

	if password != req.Password {
		return handle_errors.ErrInvalidCredentials
	}

	return nil
}
