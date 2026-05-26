package types

import (
	"context"

	"github.com/google/uuid"
)

type GetError string

const (
	ErrGetExample = GetError("Some get error")
)

type MessagesStore interface {
	Get(ctx context.Context) ([]ChatMessage, error)
	Add(ctx context.Context, content string) error
	Edit(ctx context.Context, id uuid.UUID, content string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChatsStore interface {
	Members(ctx context.Context, id string) ([]string, error)
}
