include config.env

build:
	docker compose --env-file ./config.env build

up:
	docker compose --env-file ./config.env up

down:
	docker compose --env-file ./config.env down

restart: down up

test:
	@go test -v ./...

cover:
	@go test -v -coverprofile ./tests/cover.out ./...
	@go tool cover -html ./tests/cover.out -o ./tests/cover.html
	
clean:
	docker compose --env-file ./config.env down -v --rmi all