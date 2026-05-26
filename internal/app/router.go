package app

import (
	"app/internal/app/controllers"
	"app/internal/types"
	"app/internal/ws"

	"github.com/labstack/echo/v5"
)

func SetupRouter(e *echo.Echo, hub *ws.Hub, msgStore types.MessagesStore) {
	upgrader := NewUpgrader(hub)

	e.Any("/upgrade", upgrader.Upgrade)

	chat := controllers.NewChat(msgStore)

	e.GET("/messages", chat.Load)
	e.PATCH("/messages/:id", chat.EditMessage)
	e.DELETE("/messages/:id", chat.DeleteMessage)
}
