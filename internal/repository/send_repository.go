package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SendRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewSendRepository(db *pgxpool.Pool) *SendRepositoryImpl {
	return &SendRepositoryImpl{DB: db}
}

func (r *SendRepositoryImpl) GetUserIDAndBalance(username string) (int, int32, error) {
	var userID int
	var balance int32

	err := r.DB.QueryRow(context.Background(), `SELECT id, balance FROM users WHERE username = $1`, username).Scan(&userID, &balance)
	if err != nil {
		return 0, 0, errors.New("user not found")
	}
	return userID, balance, nil
}

func (r *SendRepositoryImpl) UpdateUserBalance(userID int, newBalance int32) error {
	_, err := r.DB.Exec(context.Background(), `UPDATE users SET balance = $1 WHERE id = $2`, newBalance, userID)
	return err
}

func (r *SendRepositoryImpl) CreateTransaction(senderID, receiverID int, amount int32) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO transactions (sender_id, receiver_id, amount)
		VALUES ($1, $2, $3)
	`, senderID, receiverID, amount)
	return err
}
