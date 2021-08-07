.PHONY: test build run lint

test:
	go test ./... -coverage -covermode=atmoic -coverprofile=coverage.out

lint:
	golangci-lint run --verbose

build:
	go build -o bin/unassignederr cmd/unassignederr.go

run:
	go run cmd/unassignederr.go