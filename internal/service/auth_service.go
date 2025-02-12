package service

import (
	"avitoTT/openapi/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Authenticate(req models.AuthRequest) (models.CurrentRequest, error)
}

type AuthServiceImpl struct{}

func (s *AuthServiceImpl) Authenticate(req models.AuthRequest) (models.CurrentRequest, error) {
	if req.Username != "admin" || req.Password != "password" {
		return models.CurrentRequest{}, errors.New("Invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	// move secret key in env //
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return models.CurrentRequest{}, err
	}

	return models.CurrentRequest{Message: tokenString}, nil
}
