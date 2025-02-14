package service

import (
	"avitoTT/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceContainer struct {
	AuthService *AuthServiceImpl
	BuyService  *BuyServiceImpl
}

func NewServiceContainer(db *pgxpool.Pool) *ServiceContainer {
	authRepo := repository.NewAuthRepository(db)
	buyRepo := repository.NewBuyRepository(db)

	return &ServiceContainer{
		AuthService: NewAuthService(*authRepo),
		BuyService:  NewBuyService(*buyRepo),
	}
}
