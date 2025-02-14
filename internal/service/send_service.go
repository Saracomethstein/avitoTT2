package service

import (
	"avitoTT/internal/repository"
)

type SendService interface {
}

type SendServiceImpl struct {
	SendRepository repository.SendRepositoryImpl
}

func NewSendService(repo repository.SendRepositoryImpl) *SendServiceImpl {
	return &SendServiceImpl{SendRepository: repo}
}
