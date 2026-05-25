package store

import (
	"app/internal/types"
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

func (self *Messages) Get(ctx context.Context) ([]types.ChatMessage, error) {
	var messages []types.ChatMessage

	err := self.db.NewSelect().Model(&messages).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (self *Messages) Add(ctx context.Context, content string) error {
	msg := types.ChatMessage{Content: content}

	_, err := self.db.NewInsert().Model(&msg).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (self *Messages) Edit(ctx context.Context, id uuid.UUID, content string) error {
	msg := types.ChatMessage{
		ID:      id,
		Content: content,
	}

	_, err := self.db.NewUpdate().Model(&msg).Column("content").WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (self *Messages) Delete(ctx context.Context, id uuid.UUID) error {
	msg := types.ChatMessage{ID: id}

	_, err := self.db.NewDelete().Model(&msg).WherePK().ForceDelete().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
