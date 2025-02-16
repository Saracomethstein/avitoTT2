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

func TestSendRepository(t *testing.T) {
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
	INSERT INTO users (username, password, balance) VALUES
	('alice', 'password_hash', 1500),
	('bob', 'password_hash', 1000),
	('charlie', 'password_hash', 500);
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO transactions (sender_id, receiver_id, amount) VALUES
	(1, 2, 100),
	(2, 3, 50),
	(1, 3, 200);
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO merchandise (name, price) VALUES
	('laptop', 1000),
	('phone', 500),
	('mouse', 100);
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO purchases (user_id, merchandise_id) VALUES
	(1, 1),
	(1, 2),
	(2, 2),
	(3, 3),
	(3, 3);
	`)
	if err != nil {
		t.Fatal(err)
	}

	pgxPool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pgxPool.Close()

	infoRepo := repository.NewInfoRepository(pgxPool)

	t.Run("GetUserIDAndBalance", func(t *testing.T) {
		userID, balance, err := infoRepo.GetUserIDAndBalance("alice")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if userID != 1 {
			t.Fatalf("expected userID 1, got %d", userID)
		}
		if balance != 1500 {
			t.Fatalf("expected balance 1500, got %d", balance)
		}
	})

	t.Run("GetReceivedCoins", func(t *testing.T) {
		receivedCoins, err := infoRepo.GetReceivedCoins(3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(receivedCoins) != 2 {
			t.Fatalf("expected 2 received transactions, got %d", len(receivedCoins))
		}
		if receivedCoins[1].FromUser != "bob" || receivedCoins[1].Amount != 50 {
			t.Fatalf("expected (bob, 50), got (%s, %d)", receivedCoins[0].FromUser, receivedCoins[0].Amount)
		}
		if receivedCoins[0].FromUser != "alice" || receivedCoins[0].Amount != 200 {
			t.Fatalf("expected (alice, 200), got (%s, %d)", receivedCoins[1].FromUser, receivedCoins[1].Amount)
		}
	})

	t.Run("GetSentCoins", func(t *testing.T) {
		sentCoins, err := infoRepo.GetSentCoins(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(sentCoins) != 2 {
			t.Fatalf("expected 2 sent transactions, got %d", len(sentCoins))
		}
		if sentCoins[0].ToUser != "bob" || sentCoins[0].Amount != 100 {
			t.Fatalf("expected (bob, 100), got (%s, %d)", sentCoins[0].ToUser, sentCoins[0].Amount)
		}
		if sentCoins[1].ToUser != "charlie" || sentCoins[1].Amount != 200 {
			t.Fatalf("expected (charlie, 200), got (%s, %d)", sentCoins[1].ToUser, sentCoins[1].Amount)
		}
	})

	t.Run("GetUserInventory", func(t *testing.T) {
		inventory, err := infoRepo.GetUserInventory(3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(inventory) != 1 {
			t.Fatalf("expected 1 inventory item, got %d", len(inventory))
		}
		if inventory[0].Type != "mouse" || inventory[0].Quantity != 2 {
			t.Fatalf("expected (mouse, 2), got (%s, %d)", inventory[0].Type, inventory[0].Quantity)
		}
	})
}
