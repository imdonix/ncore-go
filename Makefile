.PHONY: test lint format

test:
	go test ./...

lint:
	golangci-lint run ./...

format:
	go fmt ./...
