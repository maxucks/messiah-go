package ws

import (
	"context"
	"log"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type WSClient struct {
	id string

	conn *websocket.Conn
	hub  *Hub

	events chan Event
}

func NewClient(conn *websocket.Conn, hub *Hub, id string) *WSClient {
	return &WSClient{
		id:     id,
		conn:   conn,
		hub:    hub,
		events: make(chan Event, 10),
	}
}

func (c *WSClient) Listen(ctx context.Context) {
	meta := c.commonMeta()
	meta.Chan = c.events

	c.hub.Add(Event{
		Meta: meta,
		Type: Connected,
	})

	go c.readLoop(ctx)
	go c.writeLoop(ctx)
}

func (c *WSClient) readLoop(ctx context.Context) {
	defer c.close()

	for {
		var e Event
		if err := wsjson.Read(ctx, c.conn, &e); err != nil {
			log.Printf("[WSClient] failed to read message: %e\n", err)
			return
		}
		e.Meta = c.commonMeta()
		c.hub.Add(e)
	}
}

func (c *WSClient) writeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-c.events:
			if !ok {
				log.Println("Write end")
				return
			}
			if err := wsjson.Write(ctx, c.conn, e); err != nil {
				log.Printf("[WSClient] failed to send message: %e\n", err)
				return
			}
		}
	}

}

func (c *WSClient) close() {
	log.Println("Closed")
	c.hub.Add(Event{
		Meta: c.commonMeta(),
		Type: Disconnected,
	})
	close(c.events)
	c.conn.Close(websocket.StatusNormalClosure, "done")
}

func (c *WSClient) commonMeta() EventMeta {
	return EventMeta{
		Sender: c.id,
	}
}
