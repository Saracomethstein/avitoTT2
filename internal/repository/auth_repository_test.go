package repository_test

import (
	"avitoTT/internal/repository"
	"avitoTT/openapi/models"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestAuthenticate(t *testing.T) {
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY,username VARCHAR(50) UNIQUE NOT NULL,password VARCHAR(256) NOT NULL,balance INTEGER DEFAULT 1000 CHECK (balance >= 0))")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", "test_user", "test_password")
	if err != nil {
		t.Fatal(err)
	}

	pgxPool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pgxPool.Close()

	authRepo := repository.NewAuthRepository(pgxPool)

	authRequest := models.AuthRequest{
		Username: "test_user",
		Password: "test_password",
	}
	err = authRepo.Authenticate(authRequest)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	authRequest.Password = "wrong_password"
	err = authRepo.Authenticate(authRequest)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != models.ErrInvalidCredentials {
		t.Fatalf("expected %v, got %v", models.ErrInvalidCredentials, err)
	}

	authRequest.Username = "test_user_001"
	authRequest.Password = "test_password_001"
	err = authRepo.Authenticate(authRequest)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", "test_user_001").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 user, got %d", count)
	}

	t.Log("Authentication tests passed successfully!")
}
