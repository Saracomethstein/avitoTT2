package repository

import (
	"avitoTT/openapi/models"
	"context"
	"errors"
	"log"
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
	log.Println("Repository: Authenticate")
	var password string

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := r.DB.QueryRow(ctx, "SELECT password FROM users WHERE username = $1", req.Username).Scan(&password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.As(err, &pgx.ErrNoRows) {
			_, err = r.DB.Exec(context.Background(), `
    INSERT INTO users (username, password)
    VALUES ($1, $2)
    ON CONFLICT (username) DO NOTHING
`, req.Username, req.Password)

			if err != nil {
				log.Println("Failed insert new user")
				return models.ErrDatabaseIssue
			}
			return nil
		}
		log.Println("Failed data base")
		return models.ErrDatabaseIssue
	}

	if password != req.Password {
		return models.ErrInvalidCredentials
	}

	return nil
}
