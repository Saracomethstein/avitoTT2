package service

import (
	handle_errors "avitoTT/internal/errors"
	"avitoTT/internal/repository"
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
			return handle_errors.ErrInvalidCredentials
		}
		return handle_errors.ErrDatabaseIssue
	}

	userID, balance, err := s.BuyRepository.GetUserIDAndBalance(username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.As(err, &pgx.ErrNoRows) {
			return handle_errors.ErrInvalidCredentials
		}
		return handle_errors.ErrDatabaseIssue
	}

	if balance < price {
		return handle_errors.ErrNotEnoughCoins
	}

	err = s.BuyRepository.MakePurchase(userID, merchID, price)
	if err != nil {
		return handle_errors.ErrDatabaseIssue
	}

	return nil
}
