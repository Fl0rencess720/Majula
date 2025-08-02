.PHONY: run deploy

include configs/.env
export

run:
	go run cmd/main.go

deploy:
	docker compose up -d --build