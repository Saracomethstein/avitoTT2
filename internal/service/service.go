package service

import (
	"avitoTT/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceContainer struct {
	AuthService *AuthServiceImpl
}

func NewServiceContainer(db *pgxpool.Pool) *ServiceContainer {
	authRepo := repository.NewAuthRepository(db)

	return &ServiceContainer{
		AuthService: NewAuthService(*authRepo),
	}
}
