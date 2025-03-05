package service

import (
	"avitoTT/internal/config"
	handle_errors "avitoTT/internal/errors"
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Authenticate(req models.AuthRequest) (models.AuthResponse, error)
}

type AuthServiceImpl struct {
	AuthRepository  repository.AuthRepositoryImpl
	RedisRepository repository.RedisRepository
	JWTSecretKey    []byte
}

func NewAuthService(repo repository.AuthRepositoryImpl, redisRepo repository.RedisRepository) *AuthServiceImpl {
	config := config.New()
	return &AuthServiceImpl{
		AuthRepository:  repo,
		RedisRepository: redisRepo,
		JWTSecretKey:    []byte(config.JWTSecretKey),
	}
}

func (s *AuthServiceImpl) Authenticate(req models.AuthRequest) (models.AuthResponse, error) {
	password, err := s.AuthRepository.CheckUserExists(req)

	switch {
	case errors.Is(err, handle_errors.ErrUserNotFound):
		if err := s.AuthRepository.InsertUser(req); err != nil {
			return models.AuthResponse{}, err
		}
	case err != nil:
		return models.AuthResponse{}, handle_errors.ErrDatabaseIssue

	case password != req.Password:
		return models.AuthResponse{}, handle_errors.ErrInvalidCredentials
	}

	if token, err := s.RedisRepository.GetCachedToken(req.Username); err == nil {
		return models.AuthResponse{Token: token}, nil
	}

	tokenString, err := s.generateToken(req.Username)
	if err != nil {
		return models.AuthResponse{}, handle_errors.ErrDatabaseIssue
	}

	if err := s.RedisRepository.CacheToken(req.Username, tokenString); err != nil {
		return models.AuthResponse{}, handle_errors.ErrDatabaseIssue
	}

	return models.AuthResponse{Token: tokenString}, nil
}

func (s *AuthServiceImpl) generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(s.JWTSecretKey)
}
