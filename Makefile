.PHONY: run

include configs/.env
export

run:
	go run cmd/main.go