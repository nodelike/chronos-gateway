.PHONY: run build test clean

# Default build options
BINARY_NAME=chronos-gateway
BUILD_DIR=./bin

# Run the application
run:
	go run cmd/gateway/main.go

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/gateway/main.go

# Run all tests
test:
	go test -v ./...

# Clean up
clean:
	rm -rf $(BUILD_DIR)

# Run with race detection for development
dev:
	go run -race cmd/gateway/main.go


docker-restart:
	docker compose down
	docker compose up -d

# Load test the API
loadtest:
	hey -n 1000 -c 50 -m POST -H "X-API-Key: test-key-1" -T "application/json" -d '{"device_id": "test-device", "user_id": "test-user", "latitude": 37.773972, "longitude": -122.431297, "event_type": "location"}' http://localhost:8080/v1/location 