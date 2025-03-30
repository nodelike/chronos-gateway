FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o chronos-gateway cmd/gateway/main.go

# Create final image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/chronos-gateway /app/chronos-gateway

# Create directory for storage
RUN mkdir -p /app/storage

# Expose ports
EXPOSE 8080 50051

# Run the app
CMD ["/app/chronos-gateway"] 