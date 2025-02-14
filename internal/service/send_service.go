package service

import (
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
	"context"
	"log"
)

type SendService interface {
}

type SendServiceImpl struct {
	SendRepository repository.SendRepositoryImpl
}

func NewSendService(repo repository.SendRepositoryImpl) *SendServiceImpl {
	return &SendServiceImpl{SendRepository: repo}
}

func (s *SendServiceImpl) SendCoin(ctx context.Context, req models.SendCoinRequest, username string) error {
	log.Println("Service: SendCoin")

	if req.ToUser == username {
		return models.ErrSendHimself
	}

	senderID, senderBalance, err := s.SendRepository.GetUserIDAndBalance(username)
	if err != nil {
		return models.ErrUserNotFound
	}

	receiverID, receiverBalance, err := s.SendRepository.GetUserIDAndBalance(req.ToUser)
	if err != nil {
		return models.ErrUserNotFound
	}

	if senderBalance < req.Amount {
		return models.ErrBalance
	}

	newSenderBalance := senderBalance - req.Amount

	tx, err := s.SendRepository.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.SendRepository.UpdateUserBalance(senderID, newSenderBalance); err != nil {
		return models.ErrDatabaseIssue
	}

	if err := s.SendRepository.UpdateUserBalance(receiverID, receiverBalance+req.Amount); err != nil {
		return models.ErrDatabaseIssue
	}

	if err := s.SendRepository.CreateTransaction(senderID, receiverID, req.Amount); err != nil {
		return models.ErrDatabaseIssue
	}

	return tx.Commit(ctx)
}
