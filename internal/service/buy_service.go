package service

import (
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
	"errors"
	"log"

	"github.com/jackc/pgx"
)

type BuyService interface {
	BuyItem(userID int, items string) error
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
		if errors.Is(err, pgx.ErrNoRows) || errors.As(err, &pgx.ErrNoRows) {
			return models.ErrInvalidCredentials
		}
		return models.ErrDatabaseIssue
	}

	userID, balance, err := s.BuyRepository.GetUserIDAndBalance(username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.As(err, &pgx.ErrNoRows) {
			return models.ErrInvalidCredentials
		}
		return models.ErrDatabaseIssue
	}

	if balance < price {
		return models.ErrNotEnoughCoins
	}

	err = s.BuyRepository.MakePurchase(userID, merchID, price)
	if err != nil {
		return models.ErrDatabaseIssue
	}

	return nil
}
