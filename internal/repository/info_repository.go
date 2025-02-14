package repository

import "github.com/jackc/pgx/v4/pgxpool"

type InfoRepositoryImpl struct {
	DB *pgxpool.Pool
}

func NewInfoRepository(db *pgxpool.Pool) *InfoRepositoryImpl {
	return &InfoRepositoryImpl{DB: db}
}
