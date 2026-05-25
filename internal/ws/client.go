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

func (self *WSClient) Listen(ctx context.Context) {
	meta := self.commonMeta()
	meta.Chan = self.events

	self.hub.Add(Event{
		Meta: meta,
		Type: Connected,
	})

	go self.readLoop(ctx)
	go self.writeLoop(ctx)
}

func (self *WSClient) readLoop(ctx context.Context) {
	defer self.close()

	for {
		var e Event
		if err := wsjson.Read(ctx, self.conn, &e); err != nil {
			log.Printf("[WSClient] failed to read message: %e\n", err)
			return
		}
		e.Meta = self.commonMeta()
		self.hub.Add(e)
	}
}

func (self *WSClient) writeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-self.events:
			if !ok {
				log.Println("Write end")
				return
			}
			if err := wsjson.Write(ctx, self.conn, e); err != nil {
				log.Printf("[WSClient] failed to send message: %e\n", err)
				return
			}
		}
	}

}

func (self *WSClient) close() {
	log.Println("Closed")
	self.hub.Add(Event{
		Meta: self.commonMeta(),
		Type: Disconnected,
	})
	close(self.events)
	self.conn.Close(websocket.StatusNormalClosure, "done")
}

func (self *WSClient) commonMeta() EventMeta {
	return EventMeta{
		Sender: self.id,
	}
}
