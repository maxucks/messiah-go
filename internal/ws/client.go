package ws

import (
	"context"
	"fmt"
	"log"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type WSClient struct {
	clientId uuid.UUID
	userId   string

	conn            *websocket.Conn
	hub             *Hub
	outcomingEvents <-chan OutcomingEvent
	closeSignal     chan string
}

func NewClient(conn *websocket.Conn, hub *Hub, userId string) *WSClient {
	return &WSClient{
		clientId:    uuid.New(),
		userId:      userId,
		conn:        conn,
		hub:         hub,
		closeSignal: make(chan string),
	}
}

func (self *WSClient) Listen() {
	self.outcomingEvents = self.hub.Register(self.clientId, self.userId)

	ctx, cancel := context.WithCancel(context.Background())

	go self.closer(cancel)
	go self.reader(ctx)
	go self.writer(ctx)

	log.Printf("socket %v opened", self.userId)
}

func (self *WSClient) reader(ctx context.Context) {
	for {
		var e IncomingEvent

		if err := wsjson.Read(ctx, self.conn, &e); err != nil {
			self.closeSignal <- fmt.Sprintf("socket read failed: %s", err)
			return
		}

		self.hub.Add(self.enrich(e))
	}
}

func (self *WSClient) writer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-self.outcomingEvents:
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
	log.Printf("closing socket %v with reason: %s", self.userId, reason)
}

func (self *WSClient) close() {
	defer self.conn.Close(websocket.StatusNormalClosure, "done")

	if err := self.hub.Unregister(self.clientId, self.userId); err != nil {
		log.Printf("unregister failed: %s\n", err)
	}

	log.Printf("socket %v closed", self.userId)
}

func (self *WSClient) enrich(e IncomingEvent) IncomingEvent {
	e.Sender = Sender{
		ClientId: self.clientId,
		UserId:   self.userId,
	}
	return e
}
