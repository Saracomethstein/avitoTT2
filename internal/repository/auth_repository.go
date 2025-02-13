package repository

import (
	"avitoTT/openapi/models"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{DB: db}
}

func (r *AuthRepositoryImpl) Authenticate(req models.AuthRequest) error {
	log.Println("Repository: Authenticate")
	var hashedPassword string

	err := r.DB.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", req.Username).Scan(&hashedPassword)

	if err != nil {
		if err != nil {
			hashedPass, hashErr := hashPassword(req.Password)
			if hashErr != nil {
				return models.ErrUserCreationFailed
			}

			_, err = r.DB.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", req.Username, hashedPass)
			if err != nil {
				log.Println("Failed insert new user")
				return models.ErrDatabaseIssue
			}
			return nil
		}
		log.Println("Failed data base")
		return models.ErrDatabaseIssue
	}

	if err := comparePassword(hashedPassword, req.Password); err != nil {
		return models.ErrInvalidCredentials
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
