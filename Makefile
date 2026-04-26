.PHONY: test lint format build docker-build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ncore cmd/ncore/main.go

test:
	go test ./...

format:
	go fmt ./...

docker:
	docker build -t ncore-go-server .
