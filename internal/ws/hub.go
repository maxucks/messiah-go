package ws

import (
	"app/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type HubClient struct {
	id     uuid.UUID
	userId string
	events chan OutcomingEvent
}

func newClient(id uuid.UUID, userId string) HubClient {
	return HubClient{
		id:     id,
		userId: userId,
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
	clients map[string]map[uuid.UUID]HubClient

	mux sync.RWMutex
}

func StartHub(messages types.MessagesStore, chats types.ChatsStore) *Hub {
	hub := &Hub{
		events:  make(chan IncomingEvent, 1000),
		clients: make(map[string]map[uuid.UUID]HubClient, 0),
		mux:     sync.RWMutex{},
		store: store{
			messages: messages,
			chats:    chats,
		},
	}

	go hub.listen()

	return hub
}

func (self *Hub) Register(clientId uuid.UUID, userId string) <-chan OutcomingEvent {
	self.mux.Lock()
	defer self.mux.Unlock()

	client := newClient(clientId, userId)

	if self.clients[userId] == nil {
		self.clients[userId] = make(map[uuid.UUID]HubClient)
	}

	self.clients[userId][clientId] = client
	log.Printf("client(id=%v, user=%v) connected\n", clientId, userId)

	return client.events
}

func (self *Hub) Unregister(clientId uuid.UUID, userId string) error {
	self.mux.Lock()
	defer self.mux.Unlock()

	client, ok := self.clients[userId][clientId]
	if !ok {
		return fmt.Errorf("no client(id=%v, user=%v) found", clientId, userId)
	}
	defer close(client.events)

	delete(self.clients[userId], clientId)
	log.Printf("client(id=%v, user=%v) disconnected\n", clientId, userId)

	return nil
}

func (self *Hub) Add(event IncomingEvent) {
	self.events <- event
}

func (self *Hub) send(sender Sender, to string, event OutcomingEvent) {
	self.mux.RLock()
	defer self.mux.RUnlock()

	clients, ok := self.clients[to]
	if !ok {
		return
	}

	for clientId, client := range clients {
		if sender.ClientId != clientId {
			client.events <- event
		}
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
			log.Printf("failed to handle event: %s\n", err)
		}
	}
}

func (self *Hub) handleDirectMessage(sender Sender, payload json.RawMessage) error {
	var msg DirectMessage
	msg.Sender = sender.UserId

	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return err
	}

	// TODO: if sender = recipient - error, can't send to yourself
	// Future - chat with self

	self.send(sender, msg.Recipient, OutcomingEvent{
		// TODO: status OK
		Type:    DirectMessageEvent,
		Payload: msg,
	})

	return nil
}

func (self *Hub) handleChatMessage(sender Sender, payload json.RawMessage) error {
	var msg ChatMessage
	msg.Sender = sender.UserId

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

	for _, userId := range userIds {
		self.send(sender, userId, OutcomingEvent{
			Type:    ChatMessageEvent,
			Payload: msg,
		})
	}

	return nil
}
