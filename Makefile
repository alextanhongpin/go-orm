include .env
export

run:
	@go run examples/main.go

up:
	@docker-compose up -d

down:
	@docker-compose down
