# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the auth service
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/auth-service ./services/auth-service/cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/auth-service .

# Copy .env file if it exists
COPY .env* ./

# Expose port
EXPOSE 8081

# Run the binary
CMD ["./auth-service"]
