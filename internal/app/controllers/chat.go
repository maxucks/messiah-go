package controllers

import (
	res "app/internal/shared"
	"app/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type ChatController struct {
	store types.MessagesStore
}

func NewChat(store types.MessagesStore) *ChatController {
	return &ChatController{store: store}
}

type LoadResponse struct {
	Messages []types.ChatMessage `json:"messages"`
}

func (self *ChatController) Load(c *echo.Context) error {
	ctx := c.Request().Context()

	messages, err := self.store.Get(ctx)
	if err != nil {
		return res.InternalError(c, err.Error())
	}

	return res.Ok(c, LoadResponse{
		Messages: messages,
	})
}

func (self *ChatController) EditMessage(c *echo.Context) error {
	ctx := c.Request().Context()

	// TODO:
	var id uuid.UUID
	content := "edited"

	err := self.store.Edit(ctx, id, content)
	if err != nil {
		return res.InternalError(c, err.Error())
	}

	return res.NoContent(c)
}

func (self *ChatController) DeleteMessage(c *echo.Context) error {
	ctx := c.Request().Context()

	// TODO:
	var id uuid.UUID

	err := self.store.Delete(ctx, id)
	if err != nil {
		return res.InternalError(c, err.Error())
	}

	return res.NoContent(c)
}
