package ws

import "encoding/json"

type EventType string

const (
	ChatMessage EventType = "chat_message"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ClientEvent struct {
	Event
	ID string `json:"-"`
}
