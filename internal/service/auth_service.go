package service

import (
	"avitoTT/internal/config"
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
	"log"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
)

var UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var PasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=]+$`)

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
	log.Println("Service: Authenticate")

	if token, err := s.RedisRepository.GetCachedToken(req.Username); err == nil {
		log.Println("Токен найден в Redis, используем его")
		return models.AuthResponse{Token: token}, nil
	}

	err := s.AuthRepository.Authenticate(req)
	if err != nil {
		return models.AuthResponse{}, err
	}

	tokenString, err := generateToken(req.Username)
	if err != nil {
		return models.AuthResponse{}, models.ErrDatabaseIssue
	}

	if err := s.RedisRepository.CacheToken(req.Username, tokenString); err != nil {
		log.Printf("Ошибка кеширования токена: %v", err)
		return models.AuthResponse{}, models.ErrDatabaseIssue
	}

	return models.AuthResponse{Token: tokenString}, nil
}

func generateToken(username string) (string, error) {
	config := config.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(config.JWTSecretKey))
}
