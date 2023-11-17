test:
	@go test ./... -v

lint:
	@echo "Running linters..."
	@golangci-lint run .

build:
	@go build -o piedpiper