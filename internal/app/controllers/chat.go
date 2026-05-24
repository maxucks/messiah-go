package controllers

import (
	res "app/internal/shared"
	"app/internal/shared/types"

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

func (cc *ChatController) Load(c *echo.Context) error {
	ctx := c.Request().Context()

	messages, err := cc.store.Get(ctx)
	if err != nil {
		return res.InternalError(c, err.Error())
	}

	return res.Ok(c, LoadResponse{
		Messages: messages,
	})
}

func (cc *ChatController) EditMessage(c *echo.Context) error {
	ctx := c.Request().Context()

	// TODO:
	var id uuid.UUID
	content := "edited"

	err := cc.store.Edit(ctx, id, content)
	if err != nil {
		return res.InternalError(c, err.Error())
	}

	return res.NoContent(c)
}

func (cc *ChatController) DeleteMessage(c *echo.Context) error {
	ctx := c.Request().Context()

	// TODO:
	var id uuid.UUID

	err := cc.store.Delete(ctx, id)
	if err != nil {
		return res.InternalError(c, err.Error())
	}

	return res.NoContent(c)
}
