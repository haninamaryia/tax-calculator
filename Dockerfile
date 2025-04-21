# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o tax-calculator main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/tax-calculator .

EXPOSE 8080

CMD ["./tax-calculator"]
