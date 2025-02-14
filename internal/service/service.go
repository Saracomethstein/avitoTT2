package service

import (
	"avitoTT/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
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
