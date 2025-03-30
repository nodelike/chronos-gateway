#!/bin/bash

# Default values
API_KEY="test-key-1"
HOST="http://localhost:8080"
DEVICE_ID="test-device-$(date +%s)"
USER_ID="test-user-123"

# Generate random coordinates around San Francisco
LAT=$(echo "37.7749 + (0.1 * $RANDOM / 32767)" | bc -l)
LNG=$(echo "-122.4194 + (0.1 * $RANDOM / 32767)" | bc -l)

# Create payload
PAYLOAD=$(cat <<EOF
{
  "device_id": "$DEVICE_ID",
  "user_id": "$USER_ID",
  "latitude": $LAT,
  "longitude": $LNG,
  "altitude": 10.0,
  "speed": 1.2,
  "heading": 270.0,
  "accuracy": 5.0,
  "event_type": "location",
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "activity_type": "walking"
}
EOF
)

echo "Sending location to $HOST/v1/location:"
echo "$PAYLOAD" | jq

# Send the request
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d "$PAYLOAD" \
  "$HOST/v1/location" | jq

echo -e "\nSent location data!"
echo "To see your event data, check the server logs." 