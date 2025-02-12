package repository

import (
	"avitoTT/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func SetupDB() *sql.DB {
	config := config.New()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)

	var db *sql.DB
	var err error
	for i := 0; i < config.DBConnectionRetries; i++ {
		db, err = sql.Open("postgres", psqlInfo)

		if err == nil {
			err = db.Ping()

			if err == nil {
				log.Println("Successfully connected to the database.")
				return db
			}
		}

		log.Printf("Retrying to connect to the database (%d/%d): %v", i+1, config.DBConnectionRetries, err)
		time.Sleep(time.Duration(config.DBConnectionDelay) * time.Second)
	}

	log.Fatalf("Failed to connect to the database after %d retries: %v", config.DBConnectionRetries, err)
	return nil
}
