package ws

import "encoding/json"

type EventType string

const (
	ChatMessageEvent   EventType = "chat_message"
	DirectMessageEvent EventType = "direct_message"
)

type IncomingEvent struct {
	Sender  string          `json:"-"`
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OutcomingEvent struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}

type DirectMessage struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

type ChatMessage struct {
	Sender string `json:"sender,omitempty"`
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}
