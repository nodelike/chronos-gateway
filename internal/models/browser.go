package models

import (
	"encoding/json"
	"time"
)

type BrowserEvent struct {
	DeviceID  string            `json:"device_id"`
	UserID    string            `json:"user_id"`
	EventType string            `json:"event_type"`
	EventData map[string]string `json:"event_data"`
	Timestamp time.Time         `json:"timestamp"`
	// Browser specific fields
	Browser    string `json:"browser"`
	BrowserVer string `json:"browser_version"`
	UserAgent  string `json:"user_agent"`
	HasMedia   bool   `json:"has_media"`
	Media      []byte `json:"-"` // Not serialized to JSON
	MediaType  string `json:"media_type,omitempty"`
}

func (e *BrowserEvent) ToJSON() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		return []byte{}
	}
	return data
}
