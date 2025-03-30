#!/bin/bash

# Default values
API_KEY="test-key-1"
HOST="http://localhost:8080"
DEVICE_ID="test-device-$(date +%s)"
USER_ID="test-user-123"
BATCH_SIZE=5

# Generate a batch of location events
PAYLOAD="["

# Start coordinates (San Francisco)
BASE_LAT=37.7749
BASE_LNG=-122.4194

for i in $(seq 1 $BATCH_SIZE); do
  # Generate coordinates with slight variations to simulate movement
  LAT=$(echo "$BASE_LAT + (0.001 * $i)" | bc -l)
  LNG=$(echo "$BASE_LNG + (0.001 * $i)" | bc -l)
  
  # Create timestamp with decreasing time (older to newer)
  TIME=$(date -u -v-${i}M +"%Y-%m-%dT%H:%M:%SZ")
  
  # Add comma for all but the last item
  if [ $i -gt 1 ]; then
    PAYLOAD="$PAYLOAD,"
  fi
  
  # Add event to batch
  PAYLOAD="$PAYLOAD
  {
    \"device_id\": \"$DEVICE_ID\",
    \"user_id\": \"$USER_ID\",
    \"latitude\": $LAT,
    \"longitude\": $LNG,
    \"altitude\": $(echo "10.0 + ($i * 0.5)" | bc -l),
    \"speed\": $(echo "1.0 + ($i * 0.2)" | bc -l),
    \"heading\": $(echo "270.0 + ($i * 5.0)" | bc -l),
    \"accuracy\": 5.0,
    \"event_type\": \"location\",
    \"timestamp\": \"$TIME\",
    \"activity_type\": \"walking\"
  }"
done

PAYLOAD="$PAYLOAD
]"

echo "Sending batch of $BATCH_SIZE locations to $HOST/v1/locations/batch:"
echo "$PAYLOAD" | jq

# Send the request
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d "$PAYLOAD" \
  "$HOST/v1/locations/batch" | jq

echo -e "\nSent batch of $BATCH_SIZE location events!"
echo "To see your event data, check the server logs." 