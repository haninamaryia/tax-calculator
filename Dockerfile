# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . ./

# Build the binary
RUN go build -o tax-calculator main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/tax-calculator .

# Copy the config.toml file from the .tests folder
COPY .tests/config.toml /app/config.toml

# Expose the port your app runs on
EXPOSE 8080

# Set the command to run the binary
CMD ["./tax-calculator"]
