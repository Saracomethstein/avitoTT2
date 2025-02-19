package service

import (
	handle_errors "avitoTT/internal/errors"
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
	"context"
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

	if req.ToUser == username {
		return handle_errors.ErrSendHimself
	}

	senderID, senderBalance, err := s.SendRepository.GetUserIDAndBalance(username)
	if err != nil {
		return handle_errors.ErrUserNotFound
	}

	receiverID, receiverBalance, err := s.SendRepository.GetUserIDAndBalance(req.ToUser)
	if err != nil {
		return handle_errors.ErrUserNotFound
	}

	if senderBalance < req.Amount {
		return handle_errors.ErrBalance
	}

	newSenderBalance := senderBalance - req.Amount

	tx, err := s.SendRepository.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.SendRepository.UpdateUserBalance(senderID, newSenderBalance); err != nil {
		return handle_errors.ErrDatabaseIssue
	}

	if err := s.SendRepository.UpdateUserBalance(receiverID, receiverBalance+req.Amount); err != nil {
		return handle_errors.ErrDatabaseIssue
	}

	if err := s.SendRepository.CreateTransaction(senderID, receiverID, req.Amount); err != nil {
		return handle_errors.ErrDatabaseIssue
	}

	return tx.Commit(ctx)
}
