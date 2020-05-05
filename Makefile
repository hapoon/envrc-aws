.PHONY: run test build

run:
	go run cmd/envrc-aws/main.go

test:
	go test ./...

build:
	go build -o bin/envrc-aws cmd/envrc-aws/main.go
