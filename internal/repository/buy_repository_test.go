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

func TestBuyRepository(t *testing.T) {
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
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(256) NOT NULL,
		balance INTEGER DEFAULT 1000 CHECK (balance >= 0)
	);
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		sender_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		receiver_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		amount INTEGER NOT NULL CHECK (amount > 0),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS merchandise (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price INTEGER NOT NULL CHECK (price > 0)
	);
	CREATE TABLE IF NOT EXISTS purchases (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		merchandise_id INTEGER REFERENCES merchandise(id) ON DELETE CASCADE,
		purchased_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO users (username, password) VALUES
	('test_user', 'password01'),
	('test_user2', 'password02');
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO merchandise (name, price) VALUES
	('t-shirt', 80),
	('cup', 20),
	('book', 50);
	`)
	if err != nil {
		t.Fatal(err)
	}

	pgxPool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pgxPool.Close()

	buyRepo := repository.NewBuyRepository(pgxPool)

	t.Run("GetMerchPrice", func(t *testing.T) {
		price, merchID, err := buyRepo.GetMerchPrice("t-shirt")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if price != 80 {
			t.Fatalf("expected price 80, got %d", price)
		}
		if merchID != 1 {
			t.Fatalf("expected merchID 1, got %d", merchID)
		}
	})

	t.Run("GetUserIDAndBalance", func(t *testing.T) {
		userID, balance, err := buyRepo.GetUserIDAndBalance("test_user")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if userID != 1 {
			t.Fatalf("expected userID 1, got %d", userID)
		}
		if balance != 1000 {
			t.Fatalf("expected balance 1000, got %d", balance)
		}
	})

	t.Run("MakePurchase", func(t *testing.T) {
		err := buyRepo.MakePurchase(1, 1, 80)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var balance int
		err = db.QueryRow("SELECT balance FROM users WHERE id = $1", 1).Scan(&balance)
		if err != nil {
			t.Fatal(err)
		}
		if balance != 920 {
			t.Fatalf("expected balance 920 after purchase, got %d", balance)
		}

		var purchaseCount int
		err = db.QueryRow("SELECT COUNT(*) FROM purchases WHERE user_id = $1 AND merchandise_id = $2", 1, 1).Scan(&purchaseCount)
		if err != nil {
			t.Fatal(err)
		}
		if purchaseCount != 1 {
			t.Fatalf("expected 1 purchase record, got %d", purchaseCount)
		}
	})

	t.Run("InsufficientBalance", func(t *testing.T) {
		err := buyRepo.MakePurchase(2, 2, 1100)
		if err == nil {
			t.Fatal("expected error due to insufficient balance, got nil")
		}
	})
}
