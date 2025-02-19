package service

import (
	handle_errors "avitoTT/internal/errors"
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Authenticate(req models.AuthRequest) (models.AuthResponse, error)
}

type AuthServiceImpl struct {
	AuthRepository  repository.AuthRepositoryImpl
	RedisRepository repository.RedisRepository
}

func NewAuthService(repo repository.AuthRepositoryImpl, redisRepo repository.RedisRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		AuthRepository:  repo,
		RedisRepository: redisRepo,
	}
}

func (s *AuthServiceImpl) Authenticate(req models.AuthRequest) (models.AuthResponse, error) {
	err := s.AuthRepository.Authenticate(req)
	if err != nil {
		return models.AuthResponse{}, err
	}

	if token, err := s.RedisRepository.GetCachedToken(req.Username); err == nil {
		return models.AuthResponse{Token: token}, nil
	}

	tokenString, err := generateToken(req.Username)
	if err != nil {
		return models.AuthResponse{}, handle_errors.ErrDatabaseIssue
	}

	if err := s.RedisRepository.CacheToken(req.Username, tokenString); err != nil {
		return models.AuthResponse{}, handle_errors.ErrDatabaseIssue
	}

	return models.AuthResponse{Token: tokenString}, nil
}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte("test_secret_key"))
}
