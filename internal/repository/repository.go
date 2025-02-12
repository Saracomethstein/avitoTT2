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

	var dbPool *pgxpool.Pool
	var err error
	for i := 0; i < config.DBConnectionRetries; i++ {
		dbPool, err = pgxpool.Connect(context.Background(), psqlInfo)

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
