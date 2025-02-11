package service

import (
	"avitoTT/openapi/models"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Authenticate(req models.AuthRequest) (models.HelloWorld, error)
}

type AuthServiceImpl struct{}

func (s *AuthServiceImpl) Authenticate(req models.AuthRequest) (models.HelloWorld, error) {
	// check user creds in Postgres //

	log.Println("data")
	log.Println(req.Username, req.Password)

	if req.Username != "admin" || req.Password != "password" {
		return models.HelloWorld{}, errors.New("Invalid credentials")
	}

	// Создаём фейковый JWT-токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	// move secret key in env //
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return models.HelloWorld{}, err
	}

	return models.HelloWorld{Message: tokenString}, nil
}
