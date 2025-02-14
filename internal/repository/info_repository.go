package repository

import (
	"avitoTT/openapi/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type InfoRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewInfoRepository(db *pgxpool.Pool) *InfoRepositoryImpl {
	return &InfoRepositoryImpl{DB: db}
}

func (r *InfoRepositoryImpl) GetUserIDAndBalance(username string) (int, int32, error) {
	var userID int
	var balance int32

	err := r.DB.QueryRow(context.Background(), `SELECT id, balance FROM users WHERE username = $1`, username).Scan(&userID, &balance)
	if err != nil {
		return 0, 0, errors.New("user not found")
	}
	return userID, balance, nil
}

func (r *InfoRepositoryImpl) GetReceivedCoins(userID int) ([]models.InfoResponseCoinHistoryReceivedInner, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT u.username, t.amount 
		FROM transactions t
		JOIN users u ON t.sender_id = u.id
		WHERE t.receiver_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receivedCoins []models.InfoResponseCoinHistoryReceivedInner
	for rows.Next() {
		var received models.InfoResponseCoinHistoryReceivedInner
		if err := rows.Scan(&received.FromUser, &received.Amount); err != nil {
			return nil, err
		}
		receivedCoins = append(receivedCoins, received)
	}
	return receivedCoins, nil
}

func (r *InfoRepositoryImpl) GetSentCoins(userID int) ([]models.InfoResponseCoinHistorySentInner, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT u.username, t.amount 
		FROM transactions t
		JOIN users u ON t.receiver_id = u.id
		WHERE t.sender_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sentCoins []models.InfoResponseCoinHistorySentInner
	for rows.Next() {
		var sent models.InfoResponseCoinHistorySentInner
		if err := rows.Scan(&sent.ToUser, &sent.Amount); err != nil {
			return nil, err
		}
		sentCoins = append(sentCoins, sent)
	}
	return sentCoins, nil
}

func (r *InfoRepositoryImpl) GetUserInventory(userID int) ([]models.InfoResponseInventoryInner, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT m.name, COUNT(p.merchandise_id)
		FROM purchases p
		JOIN merchandise m ON p.merchandise_id = m.id
		WHERE p.user_id = $1
		GROUP BY m.name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []models.InfoResponseInventoryInner
	for rows.Next() {
		var item models.InfoResponseInventoryInner
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}
	return inventory, nil
}
