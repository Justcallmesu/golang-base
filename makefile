.PHONY: dev prod build seed migrate

dev:
	air ./cmd/server/main.go

prod:
	./bin/server

build:
	go build -o bin/server ./cmd/server/main.go

seed:
	go run ./cmd/seed/main.go

migrate:
	go run ./cmd/migrate/main.go
	
