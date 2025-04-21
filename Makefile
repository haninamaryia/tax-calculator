APP_NAME=tax-calculator

# Run the app
run:
	go run main.go

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint (if using golangci-lint or similar)
lint:
	golangci-lint run

# Build binary
build:
	go build -o $(APP_NAME) main.go

# Clean build artifacts
clean:
	rm -f $(APP_NAME)

.PHONY: run test fmt lint build clean
