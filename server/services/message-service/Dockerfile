# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the api gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api-gateway ./services/api-gateway/cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/api-gateway .

# Copy .env file if it exists
COPY .env* ./

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]
