package ws

import (
	"context"
	"fmt"
	"log"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type WSClient struct {
	id string

	conn           *websocket.Conn
	hub            *Hub
	incomingEvents <-chan Event
	closeSignal    chan string
}

func NewClient(conn *websocket.Conn, hub *Hub, id string) *WSClient {
	return &WSClient{
		id:          id,
		conn:        conn,
		hub:         hub,
		closeSignal: make(chan string),
	}
}

func (self *WSClient) Listen() {
	self.incomingEvents = self.hub.Register(self.id)

	ctx, cancel := context.WithCancel(context.Background())

	go self.closer(cancel)
	go self.reader(ctx)
	go self.writer(ctx)

	log.Printf("socket %v opened", self.id)
}

func (self *WSClient) reader(ctx context.Context) {
	for {
		var e Event

		if err := wsjson.Read(ctx, self.conn, &e); err != nil {
			self.closeSignal <- fmt.Sprintf("socket read failed: %s", err)
			return
		}

		self.hub.Add(self.wrap(e))
	}
}

func (self *WSClient) writer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-self.incomingEvents:
			if !ok {
				self.closeSignal <- "incoming events chan closed"
				return
			}

			if err := wsjson.Write(ctx, self.conn, e); err != nil {
				self.closeSignal <- fmt.Sprintf("socket write failed: %s", err)
				return
			}
		}
	}
}

func (self *WSClient) closer(cancel func()) {
	defer self.close()
	defer cancel()

	reason := <-self.closeSignal
	log.Printf("closing socket %v with reason: %s", self.id, reason)
}

func (self *WSClient) close() {
	defer self.conn.Close(websocket.StatusNormalClosure, "done")

	if err := self.hub.Unregister(self.id); err != nil {
		log.Printf("unregister failed: %s\n", err)
	}

	log.Printf("socket %v closed", self.id)
}

func (self *WSClient) wrap(e Event) ClientEvent {
	return ClientEvent{
		ID:    self.id,
		Event: e,
	}
}
