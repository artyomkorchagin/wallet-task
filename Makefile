include config.env

build:
	docker compose --env-file ./config.env build

up:
	docker compose --env-file ./config.env up

down:
	docker compose --env-file ./config.env down

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
	go test -v -coverprofile ./tests/cover.out ./...
	go tool cover -html ./tests/cover.out -o ./tests/cover.html
	./tests/cover.html
	
clean:
	docker compose --env-file ./config.env down -v --rmi all