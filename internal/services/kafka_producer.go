package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer        sarama.AsyncProducer
	topicMap        map[string]string // source -> topic mapping
	developmentMode bool
}

func NewKafkaProducer(brokers []string, developmentMode bool) *KafkaProducer {
	// Create topic mapping
	topicMap := map[string]string{
		"android":  "android-events",
		"macos":    "macos-events",
		"browser":  "browser-events",
		"location": "location-events",
	}

	// If in development mode, we don't connect to Kafka
	if developmentMode {
		log.Println("Starting Kafka producer in development mode - messages will be logged, not sent to Kafka")
		return &KafkaProducer{
			producer:        nil,
			topicMap:        topicMap,
			developmentMode: true,
		}
	}

	// Configure Kafka producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.ClientID = "chronos-gateway"

	// Add retry logic with a shorter timeout
	config.Net.DialTimeout = 5 * time.Second
	config.Net.ReadTimeout = 5 * time.Second
	config.Net.WriteTimeout = 5 * time.Second

	// Try each broker until one works
	var connectedBrokers []string
	var producer sarama.AsyncProducer
	var err error

	log.Printf("[KAFKA] Attempting to connect to brokers: %v", brokers)

	// First try all brokers together
	producer, err = sarama.NewAsyncProducer(brokers, config)
	if err == nil {
		log.Printf("[KAFKA] Successfully connected to brokers: %v", brokers)
	} else {
		log.Printf("[KAFKA] Error connecting to all brokers together: %v", err)
		log.Println("[KAFKA] Will try each broker individually")

		// Try each broker individually
		for _, broker := range brokers {
			singleBroker := []string{broker}
			producer, err = sarama.NewAsyncProducer(singleBroker, config)
			if err == nil {
				log.Printf("[KAFKA] Successfully connected to broker: %s", broker)
				connectedBrokers = append(connectedBrokers, broker)
				break
			} else {
				log.Printf("[KAFKA] Failed to connect to broker %s: %v", broker, err)
			}
		}
	}

	// If all connection attempts failed, fall back to development mode
	if err != nil || producer == nil {
		log.Printf("[KAFKA] Error creating Kafka producer: %v", err)
		log.Println("[KAFKA] Falling back to development mode - messages will be logged")
		return &KafkaProducer{
			producer:        nil,
			topicMap:        topicMap,
			developmentMode: true,
		}
	}

	// Start a goroutine to handle success and error messages
	go func() {
		for {
			select {
			case success := <-producer.Successes():
				log.Printf("[KAFKA] Successfully sent message to topic %s partition %d offset %d",
					success.Topic, success.Partition, success.Offset)
			case err := <-producer.Errors():
				log.Printf("[KAFKA] Failed to send message: %v", err)
			}
		}
	}()

	return &KafkaProducer{
		producer:        producer,
		topicMap:        topicMap,
		developmentMode: false,
	}
}

func (kp *KafkaProducer) SendEvent(source string, event []byte) {
	topic, exists := kp.topicMap[source]
	if !exists {
		// Default to source name as topic if not in map
		topic = source + "-events"
	}

	// In development mode, just log the message
	if kp.developmentMode {
		// Pretty print JSON if possible
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal(event, &prettyJSON); err == nil {
			prettyJSONBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
			log.Printf("[DEV MODE] Would send to topic %s:\n%s", topic, string(prettyJSONBytes))
		} else {
			// Fallback to raw bytes if not valid JSON
			log.Printf("[DEV MODE] Would send to topic %s: %s", topic, string(event))
		}
		return
	}

	// Send to Kafka in production mode
	log.Printf("[KAFKA] Sending message to topic %s", topic)
	kp.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(event),
	}
}
