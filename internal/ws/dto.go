package ws

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventType string

const (
	ChatMessageEvent   EventType = "chat.message"
	DirectMessageEvent EventType = "direct.message"
)

type Sender struct {
	UserId   string    `json:"-"`
	ClientId uuid.UUID `json:"-"`
}

type IncomingEvent struct {
	Sender  Sender          `json:"-"`
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
	ChatID string `json:"chatId"`
	Text   string `json:"text"`
}
