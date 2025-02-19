package service

import (
	handle_errors "avitoTT/internal/errors"
	"avitoTT/internal/repository"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type ServiceContainer struct {
	AuthService *AuthServiceImpl
	BuyService  *BuyServiceImpl
	InfoService *InfoServiceImpl
	SendService *SendServiceImpl
}

func NewServiceContainer(db *pgxpool.Pool) *ServiceContainer {
	authRepo := repository.NewAuthRepository(db)
	buyRepo := repository.NewBuyRepository(db)
	infoRepo := repository.NewInfoRepository(db)
	sendRepo := repository.NewSendRepository(db)
	redisRepo := repository.NewRedisRepository()

	return &ServiceContainer{
		AuthService: NewAuthService(*authRepo, *redisRepo),
		BuyService:  NewBuyService(*buyRepo),
		InfoService: NewInfoService(*infoRepo),
		SendService: NewSendService(*sendRepo),
	}
}

func ExtractTokenFromHeader(ctx echo.Context) (string, error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", handle_errors.ErrMissingToken
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", handle_errors.ErrInvalidTokenFormat
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func ExtractUsernameFromToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("test_secret_key"), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", handle_errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", handle_errors.ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", handle_errors.ErrInvalidToken
	}

	return username, nil
}
