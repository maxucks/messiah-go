package ws

import (
	"app/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type HubClient struct {
	id     string
	events chan OutcomingEvent
}

func newClient(id string) HubClient {
	return HubClient{
		id:     id,
		events: make(chan OutcomingEvent, 1000),
	}
}

type store struct {
	messages types.MessagesStore
	chats    types.ChatsStore
}

type Hub struct {
	store store

	events  chan IncomingEvent
	clients map[string]HubClient

	mux sync.RWMutex
}

func StartHub(messages types.MessagesStore, chats types.ChatsStore) *Hub {
	hub := &Hub{
		events:  make(chan IncomingEvent, 1000),
		clients: make(map[string]HubClient, 0),
		mux:     sync.RWMutex{},
		store: store{
			messages: messages,
			chats:    chats,
		},
	}

	go hub.listen()

	return hub
}

func (self *Hub) Register(id string) <-chan OutcomingEvent {
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

func (self *Hub) Add(e IncomingEvent) {
	self.events <- e
}

func (self *Hub) send(id string, e OutcomingEvent) {
	self.mux.RLock()
	defer self.mux.RUnlock()

	client, ok := self.clients[id]
	if ok {
		client.events <- e
	}
}

func (self *Hub) listen() {
	for e := range self.events {
		var err error

		switch e.Type {
		case ChatMessageEvent:
			err = self.handleChatMessage(e.Sender, e.Payload)
		case DirectMessageEvent:
			err = self.handleDirectMessage(e.Sender, e.Payload)
		default:
			fmt.Printf("client %v sent unkown event '%v': %v\n", e.Sender, e.Type, e.Payload)
		}

		if err != nil {
			log.Println(err)
		}
	}
}

func (self *Hub) handleDirectMessage(sender string, payload json.RawMessage) error {
	var msg DirectMessage
	msg.Sender = sender

	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return err
	}

	// TODO: if sender = recipient - error, can't send to yourself
	// Future - chat with self

	self.send(msg.Recipient, OutcomingEvent{
		// TODO: status OK
		Type:    DirectMessageEvent,
		Payload: msg,
	})

	return nil
}

func (self *Hub) handleChatMessage(sender string, payload json.RawMessage) error {
	var msg ChatMessage
	msg.Sender = sender

	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// TODO: Cache it
	userIds, err := self.store.chats.Members(ctx, msg.ChatID)
	if err != nil {
		return err
	}

	for _, id := range userIds {
		if id == sender {
			continue
		}

		self.send(id, OutcomingEvent{
			Type:    ChatMessageEvent,
			Payload: msg,
		})
	}

	return nil
}
