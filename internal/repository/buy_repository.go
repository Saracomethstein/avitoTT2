package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type BuyRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewBuyRepository(db *pgxpool.Pool) *BuyRepositoryImpl {
	return &BuyRepositoryImpl{DB: db}
}

func (r *BuyRepositoryImpl) GetMerchPrice(itemName string) (int, int, error) {
	var price int
	var merchID int
	err := r.DB.QueryRow(context.Background(),
		"SELECT price, id FROM merchandise WHERE name = $1", itemName).Scan(&price, &merchID)

	return price, merchID, err
}

func (r *BuyRepositoryImpl) GetUserIDAndBalance(username string) (int, int, error) {
	var userID int
	var balance int
	err := r.DB.QueryRow(context.Background(),
		"SELECT id, balance FROM users WHERE username = $1", username).Scan(&userID, &balance)

	return userID, balance, err
}

func (r *BuyRepositoryImpl) MakePurchase(userID, merchID, price int) error {
	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE id = $2", price, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO purchases (user_id, merchandise_id) VALUES ($1, $2)",
		userID, merchID)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
