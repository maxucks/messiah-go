package ws

import "encoding/json"

type EventType string

const (
	Connected    EventType = "connected"
	Disconnected EventType = "disconnected"
	ChatMessage  EventType = "chat_message"
)

type EventMeta struct {
	Sender string
	Chan   chan<- Event
}

type Event struct {
	Meta    EventMeta       `json:"-"`
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
