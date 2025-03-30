package services

import (
	"encoding/json"
	"log"

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
	producer, err := sarama.NewAsyncProducer(brokers, config)

	if err != nil {
		log.Printf("Error creating Kafka producer: %v", err)
		log.Println("Falling back to development mode - messages will be logged")
		return &KafkaProducer{
			producer:        nil,
			topicMap:        topicMap,
			developmentMode: true,
		}
	}

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
	kp.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(event),
	}
}
