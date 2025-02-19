package service

import (
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
)

type InfoService interface {
}

type InfoServiceImpl struct {
	InfoRepository repository.InfoRepositoryImpl
}

func NewInfoService(repo repository.InfoRepositoryImpl) *InfoServiceImpl {
	return &InfoServiceImpl{InfoRepository: repo}
}

func (s *InfoServiceImpl) GetUserInfo(username string) (models.InfoResponse, error) {
	var response models.InfoResponse

	userID, balance, err := s.InfoRepository.GetUserIDAndBalance(username)
	if err != nil {
		return models.InfoResponse{}, err
	}
	response.Coins = balance

	receivedCoins, err := s.InfoRepository.GetReceivedCoins(userID)
	if err != nil {
		return models.InfoResponse{}, err
	}
	response.CoinHistory.Received = receivedCoins

	sentCoins, err := s.InfoRepository.GetSentCoins(userID)
	if err != nil {
		return models.InfoResponse{}, err
	}
	response.CoinHistory.Sent = sentCoins

	inventory, err := s.InfoRepository.GetUserInventory(userID)
	if err != nil {
		return models.InfoResponse{}, err
	}
	response.Inventory = inventory

	return response, nil
}
