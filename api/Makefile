BIN_DIR=bin

build:
	mkdir -p $(BIN_DIR) #
	go build -o $(BIN_DIR)/ -v ./cmd/app

build-migrator:
	mkdir -p $(BIN_DIR) #
	go build -o $(BIN_DIR)/ -v ./cmd/migrator

migrate.up:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/time_tracker?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/time_tracker?sslmode=disable' down