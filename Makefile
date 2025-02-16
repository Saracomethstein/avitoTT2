BINARY_NAME=build/avito_app

.PHONY: all build run clean docker-build docker-run docker-clean

all: run

deps:
	@echo "==> Installing dependencies..."
	go mod tidy

build: deps
	@echo "==> Building the application..."
	mkdir build
	go build -o $(BINARY_NAME) cmd/avitoTT/main.go

run: build
	@echo "==> Running the application..."
	./$(BINARY_NAME)

clean:
	@echo "==> Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)
	rm -rf build

docker-build:
	@echo "==> Building Docker containers..."
	docker compose build

docker-up: docker-build
	@echo "==> Starting Docker containers..."
	docker compose up

docker-down:
	@echo "==> Stopping Docker containers..."
	docker compose down

test:
	go test internal/repository/auth_repository_test.go
	go test internal/repository/send_repository_test.go
	go test internal/repository/buy_repository_test.go
	go test internal/repository/info_repository_test.go
	go test internal/errors/errors_test.go

postman_test:
	newman run tests/postman_collection_test/AvitoTT\ \(auth\).postman_collection.json 
	newman run tests/postman_collection_test/AvitoTT\ \(send\ coins\).postman_collection.json 
	newman run tests/postman_collection_test/AvitoTT\ \(info\).postman_collection.json   
	newman run tests/postman_collection_test/AvitoTT\ \(buy\ item\).postman_collection.json   

load_test:
	k6 run tests/k6_test/test.js