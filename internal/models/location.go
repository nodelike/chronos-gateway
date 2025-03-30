package models

import (
	"encoding/json"
	"time"
)

// LocationEvent represents geolocation data from React Native Background Geolocation
type LocationEvent struct {
	// Core location data
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"device_id"`
	UserID    string    `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude,omitempty"`
	Speed     float64   `json:"speed,omitempty"`
	Heading   float64   `json:"heading,omitempty"`
	Accuracy  float64   `json:"accuracy,omitempty"`
	EventType string    `json:"event_type"` // "location", "motion", "geofence", etc.
	
	GeofenceID     string `json:"geofence_id,omitempty"`
	ActivityType string  `json:"activity_type,omitempty"` // "still", "walking", "in_vehicle", etc.

}

func (e *LocationEvent) ToJSON() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		return []byte{}
	}
	return data
}
