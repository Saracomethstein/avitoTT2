package service

import (
	"avitoTT/internal/config"
	"avitoTT/internal/repository"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
)

type BuyService interface {
	BuyItem(userID int, items string) error
	ExtractUserIDFromToken(ctx echo.Context) (int, error)
}

type BuyServiceImpl struct {
	BuyRepository repository.BuyRepositoryImpl
}

func NewBuyService(repo repository.BuyRepositoryImpl) *BuyServiceImpl {
	return &BuyServiceImpl{BuyRepository: repo}
}

func (s *BuyServiceImpl) BuyItem(username, items string) error {
	log.Println("Service: BuyItem")

	price, merchID, err := s.BuyRepository.GetMerchPrice(items)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("merch item not found")
		}
		return err
	}

	userID, err := s.BuyRepository.GetUserID(username)
	if err != nil {
		return err
	}

	balance, err := s.BuyRepository.GetUserBalance(userID)
	if err != nil {
		return err
	}

	if balance < price {
		return fmt.Errorf("not enough coins to buy %s", items)
	}

	err = s.BuyRepository.MakePurchase(userID, merchID, price)
	if err != nil {
		return err
	}

	return nil
}

func (c *BuyServiceImpl) ExtractTokenFromHeader(ctx echo.Context) (string, error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func (s *BuyServiceImpl) ExtractUsernameFromToken(tokenStr string) (string, error) {
	config := config.New()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("username not found in token")
	}

	return username, nil
}
