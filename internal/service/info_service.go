package service

import "avitoTT/internal/repository"

type InfoService interface {
}

type InfoServiceImpl struct {
	InfoRepository repository.InfoRepositoryImpl
}

func NewInfoService(repo repository.InfoRepositoryImpl) *InfoServiceImpl {
	return &InfoServiceImpl{InfoRepository: repo}
}
