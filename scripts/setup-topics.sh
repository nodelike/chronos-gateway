#!/bin/bash
set -e

echo "Starting Kafka topic setup..."

# Function to check if Kafka is ready
check_kafka() {
  echo "Testing connection to Kafka..."
  if kafka-broker-api-versions.sh --bootstrap-server kafka:9092 > /dev/null 2>&1; then
    echo "Successfully connected to Kafka"
    return 0
  else
    echo "Cannot connect to Kafka yet"
    return 1
  fi
}

# Wait until Kafka is ready with a timeout
MAX_ATTEMPTS=10
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
  echo "Attempt $((ATTEMPT+1))/$MAX_ATTEMPTS: Checking if Kafka is ready..."
  
  if check_kafka; then
    # Create the topic
    echo "Creating topic: location-events"
    kafka-topics.sh --bootstrap-server kafka:9092 --create --if-not-exists --topic location-events --partitions 3 --replication-factor 1
    
    # Verify the topic was created
    echo "Verifying topic creation:"
    kafka-topics.sh --bootstrap-server kafka:9092 --list
    
    echo "Kafka topics created successfully"
    exit 0
  fi
  
  ATTEMPT=$((ATTEMPT+1))
  echo "Waiting 5 seconds before next attempt..."
  sleep 5
done

echo "Failed to connect to Kafka after $MAX_ATTEMPTS attempts"
exit 1 