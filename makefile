.PHONY: dev prod build seed migrate test cover lint run

dev:
	air ./cmd/server/main.go

run: build
	./bin/server

prod:
	./bin/server

test:
	go test  ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build:
	go build -o bin/server ./cmd/server

lint:
	golangci-lint run ./...

seed:
	go run ./cmd/seeder/main.go

migrate:
	go run ./cmd/migrations/main.go
	
