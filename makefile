.PHONY: dev prod build seed migrate

dev:
	air ./cmd/server/main.go

prod:
	./bin/server

build:
	go build -o bin/server ./cmd/server/main.go && ./bin/server

seed:
	go run ./cmd/seeder/main.go

migrate:
	go run ./cmd/migrations/main.go
	
