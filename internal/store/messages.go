package store

import (
	"app/internal/shared/types"
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Messages struct {
	db *bun.DB
}

func NewMessages(db *bun.DB) *Messages {
	return &Messages{db}
}

func (e *Messages) Get(ctx context.Context) ([]types.ChatMessage, error) {
	var messages []types.ChatMessage

	err := e.db.NewSelect().Model(&messages).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (e *Messages) Add(ctx context.Context, content string) error {
	msg := types.ChatMessage{Content: content}

	_, err := e.db.NewInsert().Model(&msg).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (e *Messages) Edit(ctx context.Context, id uuid.UUID, content string) error {
	msg := types.ChatMessage{
		ID:      id,
		Content: content,
	}

	_, err := e.db.NewUpdate().Model(&msg).Column("content").WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (e *Messages) Delete(ctx context.Context, id uuid.UUID) error {
	msg := types.ChatMessage{ID: id}

	_, err := e.db.NewDelete().Model(&msg).WherePK().ForceDelete().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
