package ws

import (
	"app/internal/types"
	"context"
	"fmt"
)

type HubClient struct {
	events chan<- Event
}

type Hub struct {
	store types.MessagesStore

	events  chan Event
	clients map[string]HubClient
}

func StartHub(store types.MessagesStore) *Hub {
	hub := &Hub{
		events:  make(chan Event, 1000),
		clients: make(map[string]HubClient, 0),
		store:   store,
	}

	go hub.listen()

	return hub
}

func (self *Hub) Add(e Event) {
	self.events <- e
}

func (self *Hub) send(userId string, e Event) {
	u, ok := self.clients[userId]
	if ok {
		// FIXME: could be potentially closed
		u.events <- e
	}
}

func (self *Hub) listen() {
	for e := range self.events {
		switch e.Type {
		case Connected:
			self.clients[e.Meta.Sender] = HubClient{events: e.Meta.Chan}
			fmt.Printf("client %v connected\n", e.Meta.Sender)
		case Disconnected:
			delete(self.clients, e.Meta.Sender)
			fmt.Printf("client %v disconnected\n", e.Meta.Sender)
		case ChatMessage:
			self.store.Add(context.TODO(), string(e.Payload))
			self.send(e.Meta.Sender, e)
		default:
			fmt.Printf("client %v sent unkown event '%v': %v\n", e.Meta.Sender, e.Type, e.Payload)
		}
	}
}
