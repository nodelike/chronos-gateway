package models

import (
	"encoding/json"
	"time"
)

type MacOSEvent struct {
	DeviceID    string            `json:"device_id"`
	UserID      string            `json:"user_id"`
	EventType   string            `json:"event_type"`
	EventData   map[string]string `json:"event_data"`
	Timestamp   time.Time         `json:"timestamp"`
	AppVersion  string            `json:"app_version"`
	OSVersion   string            `json:"os_version"`
	DeviceModel string            `json:"device_model"`
	// MacOS specific fields
	DesktopEnv string `json:"desktop_env"`
}

func (e *MacOSEvent) ToJSON() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		return []byte{}
	}
	return data
}
