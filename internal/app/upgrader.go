package app

import (
	res "app/internal/shared"
	"app/internal/ws"
	"context"
	"net/http"

	"github.com/coder/websocket"
	"github.com/labstack/echo/v5"
)

type WSUpgrader struct {
	hub *ws.Hub
}

func NewUpgrader(hub *ws.Hub) *WSUpgrader {
	return &WSUpgrader{hub}
}

func (self *WSUpgrader) Upgrade(c *echo.Context) error {
	userId := c.QueryParam("user_id")
	if userId == "" {
		c.Logger().Error("user_id is required")
		return c.JSON(http.StatusBadRequest, res.ErrorResponse("user_id is required"))
	}

	options := &websocket.AcceptOptions{InsecureSkipVerify: true}
	conn, err := websocket.Accept(c.Response(), c.Request(), options)
	if err != nil {
		c.Logger().Error("failed to accept ws", "error", err)
		return c.JSON(http.StatusInternalServerError, res.ErrorResponse(err.Error()))
	}

	ctx := context.Background()
	client := ws.NewClient(conn, self.hub, userId)
	client.Listen(ctx)

	return nil
}
