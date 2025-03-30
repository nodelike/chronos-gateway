package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nodelike/chronos-gateway/internal/services"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

func WebSocketHandler(producer *services.KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}
		defer conn.Close()

		// Set read deadline
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Set ping handler to keep connection alive
		conn.SetPingHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})

		// Process messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			// Process the message based on the format
			var event map[string]interface{}
			if err := json.Unmarshal(message, &event); err != nil {
				log.Printf("Error parsing JSON: %v", err)
				continue
			}

			// Determine the source type
			source, ok := event["source"].(string)
			if !ok {
				source = "unknown"
			}

			// Send to Kafka
			producer.SendEvent(source, message)

			// Send acknowledgment
			response := map[string]interface{}{
				"status":    "received",
				"timestamp": time.Now().Unix(),
			}

			responseJSON, _ := json.Marshal(response)
			if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
				log.Printf("Error sending response: %v", err)
				break
			}
		}
	}
}
