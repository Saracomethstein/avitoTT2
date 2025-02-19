package repository

import (
	"avitoTT/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupDB() *pgxpool.Pool {
	config := config.New()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)

	psqCfg, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		log.Fatalf("Unable to parse database configuration: %v", err)
	}

	psqCfg.MaxConns = 100
	psqCfg.MinConns = 10
	psqCfg.MaxConnLifetime = time.Hour

	var dbPool *pgxpool.Pool
	for i := 0; i < config.DBConnectionRetries; i++ {
		dbPool, err = pgxpool.ConnectConfig(context.Background(), psqCfg)

		if err == nil {
			err = dbPool.Ping(context.Background())

			if err == nil {
				log.Println("Successfully connected to the database.")
				return dbPool
			}
		}

		log.Printf("Retrying to connect to the database (%d/%d): %v", i+1, config.DBConnectionRetries, err)
		time.Sleep(time.Duration(config.DBConnectionDelay) * time.Second)
	}

	log.Fatalf("Failed to connect to the database after %d retries: %v", config.DBConnectionRetries, err)
	return nil
}
