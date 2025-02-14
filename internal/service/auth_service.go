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
	AuthRepository repository.AuthRepositoryImpl
}

func NewAuthService(repo repository.AuthRepositoryImpl) *AuthServiceImpl {
	return &AuthServiceImpl{AuthRepository: repo}
}

func (s *AuthServiceImpl) Authenticate(req models.AuthRequest) (models.AuthResponse, error) {
	log.Println("Service: Authenticate")

	if !UsernameRegex.MatchString(req.Username) || !PasswordRegex.MatchString(req.Password) {
		return models.AuthResponse{}, models.ErrInvalidCredentials
	}

	err := s.AuthRepository.Authenticate(req)
	if err != nil {
		return models.AuthResponse{}, err
	}

	tokenString, err := generateToken(req.Username)
	if err != nil {
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
