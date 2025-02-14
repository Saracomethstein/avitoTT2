package repository

import "github.com/jackc/pgx/v4/pgxpool"

type SendRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewSendRepository(db *pgxpool.Pool) *SendRepositoryImpl {
	return &SendRepositoryImpl{DB: db}
}
