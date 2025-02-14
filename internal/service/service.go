package service

import (
	"avitoTT/internal/config"
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
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

	return &ServiceContainer{
		AuthService: NewAuthService(*authRepo),
		BuyService:  NewBuyService(*buyRepo),
		InfoService: NewInfoService(*infoRepo),
		SendService: NewSendService(*sendRepo),
	}
}

func ExtractTokenFromHeader(ctx echo.Context) (string, error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", models.ErrMissingToken
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", models.ErrInvalidTokenFormat
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func ExtractUsernameFromToken(tokenStr string) (string, error) {
	config := config.New()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", models.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", models.ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", models.ErrInvalidToken
	}

	return username, nil
}
