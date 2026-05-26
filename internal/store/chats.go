package store

import "context"

type chats struct {
	// db *bun.DB
}

func NewChats() *chats {
	return &chats{}
}

func (self *chats) Members(ctx context.Context, id string) ([]string, error) {
	if id != "1" {
		return []string{}, nil
	}

	userIds := []string{"1", "2", "3"}
	return userIds, nil
}
