package store

import (
	"app/internal/types"
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type messages struct {
	db *bun.DB
}

func NewMessages(db *bun.DB) *messages {
	return &messages{db}
}

func (self *messages) Get(ctx context.Context) ([]types.ChatMessage, error) {
	var res []types.ChatMessage

	err := self.db.NewSelect().Model(&res).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (self *messages) Add(ctx context.Context, content string) error {
	msg := types.ChatMessage{Content: content}

	_, err := self.db.NewInsert().Model(&msg).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (self *messages) Edit(ctx context.Context, id uuid.UUID, content string) error {
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

func (self *messages) Delete(ctx context.Context, id uuid.UUID) error {
	msg := types.ChatMessage{ID: id}

	_, err := self.db.NewDelete().Model(&msg).WherePK().ForceDelete().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
