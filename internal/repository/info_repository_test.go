package repository_test

import (
	"avitoTT/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInfoRepository(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(256) NOT NULL,
		balance INTEGER DEFAULT 1000 CHECK (balance >= 0)
	);
	CREATE TABLE transactions (
		id SERIAL PRIMARY KEY,
		sender_id INTEGER REFERENCES users(id),
		receiver_id INTEGER REFERENCES users(id),
		amount INTEGER NOT NULL CHECK (amount > 0)
	);
	CREATE TABLE merchandise (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL
	);
	CREATE TABLE purchases (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		merchandise_id INTEGER REFERENCES merchandise(id)
	);
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO users (username, password, balance) VALUES
	('alice', 'pass01', 1500),
	('bob', 'pass02', 1000);
	`)
	if err != nil {
		t.Fatal(err)
	}

	pgxPool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pgxPool.Close()

	repo := repository.NewInfoRepository(pgxPool)

	t.Run("GetUserIDAndBalance", func(t *testing.T) {
		userID, balance, err := repo.GetUserIDAndBalance("alice")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if userID != 1 || balance != 1500 {
			t.Fatalf("unexpected result: userID=%d, balance=%d", userID, balance)
		}
	})
}
