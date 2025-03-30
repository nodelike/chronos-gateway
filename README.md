# Chronos Gateway

A multi-protocol API Gateway for the Chronos system that handles data collection from various client sources, including geolocation data from mobile devices.

## Features

- Multi-protocol support (HTTP, WebSocket, gRPC)
- Authentication using API keys
- Request validation and payload processing
- Metrics collection with Prometheus
- Async Kafka message production
- Media file handling with MinIO
- Location tracking for mobile applications
- Observability instrumentation

## Setup

### Prerequisites

- Go 1.24 or higher
- Kafka cluster
- MinIO server (for media uploads)
- Prometheus (optional, for metrics)

### Running the Service

1. Clone the repository
2. Configure settings in `configs/config.yaml`
3. Run the service:

```bash
go run cmd/gateway/main.go
```

Or build and run the binary:

```bash
go build -o gateway cmd/gateway/main.go
./gateway
```

## API Endpoints

### HTTP Endpoints

- `POST /v1/android` - Collect events from Android devices
- `POST /v1/macos` - Collect events from MacOS devices 
- `POST /v1/browser` - Collect events from web browsers
- `POST /v1/location` - Collect location data from mobile devices
- `POST /v1/locations/batch` - Collect batched location data from mobile devices
- `POST /v1/upload` - Upload media files

### WebSocket Endpoint

- `GET /v1/ws` - WebSocket connection for real-time events

### gRPC Service

The gRPC service is available on port 50051 (default) and supports event collection.

## Location Tracking

The gateway supports real-time location tracking from React Native apps using react-native-background-geolocation. Features include:

- Single location event capture
- Batch processing for efficient uploads
- Metrics for tracking latency and batch sizes
- Support for geofence events and activity recognition

### GPS Data Collection Flow

1. Mobile app collects GPS coordinates using react-native-background-geolocation
2. App batches multiple location points for efficient transmission
3. App sends batch to API Gateway via POST /v1/locations/batch
4. Gateway validates, authenticates, and records metrics
5. Gateway sends each location point to Kafka (location-events topic)
6. Stream processor consumes events for enrichment and storage
7. Time-series data is stored in TimescaleDB for analysis

### Example Location Event

```json
{
  "timestamp": "2023-06-15T12:30:45Z",
  "device_id": "d8e8fca2-dc0f-4d1a-a7a1-ca5dfbd6c9b1",
  "user_id": "user123",
  "latitude": 37.773972,
  "longitude": -122.431297,
  "altitude": 10.0,
  "speed": 1.2,
  "heading": 270.0,
  "accuracy": 5.0,
  "event_type": "location",
  "activity_type": "walking",
  "is_moving": true
}
```

## Metrics

The service exposes Prometheus metrics at the `/metrics` endpoint, including:

- API request counts and latencies
- Location event counts by type
- Location event latency (time between creation and reception)
- Batch size distribution

## Client Authentication

Clients must include an API key in the `X-API-Key` header for authentication. 