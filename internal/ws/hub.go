package ws

import (
	"app/internal/shared/types"
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

func (h *Hub) Add(e Event) {
	h.events <- e
}

func (h *Hub) send(userId string, e Event) {
	u, ok := h.clients[userId]
	if ok {
		// FIXME: could be potentially closed
		u.events <- e
	}
}

func (h *Hub) listen() {
	for e := range h.events {
		switch e.Type {
		case Connected:
			h.clients[e.Meta.Sender] = HubClient{events: e.Meta.Chan}
			fmt.Printf("client %v connected\n", e.Meta.Sender)
		case Disconnected:
			delete(h.clients, e.Meta.Sender)
			fmt.Printf("client %v disconnected\n", e.Meta.Sender)
		case ChatMessage:
			h.store.Add(context.TODO(), string(e.Payload))
			h.send(e.Meta.Sender, e)
		default:
			fmt.Printf("client %v sent unkown event '%v': %v\n", e.Meta.Sender, e.Type, e.Payload)
		}
	}
}
