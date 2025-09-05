include .env

build:
	docker compose build

up:
	docker compose up

down:
	docker compose down

restart: down up

db-up:
	@echo "Running migrations..."
	@goose -dir migrations postgres "$(DB_DSN)" up

db-down:
	@echo "Rolling back migrations..."
	@goose -dir migrations postgres "$(DB_DSN)" down

db-status:
	@goose -dir migrations postgres "$(DB_DSN)" status

tests:
	go test ./internal/repository/postgres/user/... -v
	
clean:
	docker compose down -v --rmi all