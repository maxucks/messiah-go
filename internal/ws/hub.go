package ws

import (
	"app/internal/types"
	"fmt"
	"log"
	"sync"
)

type HubClient struct {
	id     string
	events chan Event
}

func newClient(id string) HubClient {
	return HubClient{
		id:     id,
		events: make(chan Event, 1000),
	}
}

type Hub struct {
	store types.MessagesStore

	events  chan ClientEvent
	clients map[string]HubClient

	mux sync.RWMutex
}

func StartHub(store types.MessagesStore) *Hub {
	hub := &Hub{
		events:  make(chan ClientEvent, 1000),
		clients: make(map[string]HubClient, 0),
		store:   store,
		mux:     sync.RWMutex{},
	}

	go hub.listen()

	return hub
}

func (self *Hub) Register(id string) <-chan Event {
	self.mux.Lock()
	defer self.mux.Unlock()

	client := newClient(id)
	self.clients[id] = client
	log.Printf("client(id=%v) connected\n", id)

	return client.events
}

func (self *Hub) Unregister(id string) error {
	self.mux.Lock()
	defer self.mux.Unlock()

	client, ok := self.clients[id]
	if !ok {
		return fmt.Errorf("no client(id=%v) found", id)
	}
	defer close(client.events)

	delete(self.clients, id)
	log.Printf("client(id=%v) disconnected\n", id)

	return nil
}

func (self *Hub) Add(e ClientEvent) {
	self.events <- e
}

func (self *Hub) send(id string, e Event) {
	self.mux.RLock()
	defer self.mux.RUnlock()

	client, ok := self.clients[id]
	if ok {
		client.events <- e
	}
}

func (self *Hub) listen() {
	for e := range self.events {
		switch e.Type {
		case ChatMessage:
			self.send(e.ID, e.Event)
		default:
			fmt.Printf("client %v sent unkown event '%v': %v\n", e.ID, e.Type, e.Payload)
		}
	}
}
