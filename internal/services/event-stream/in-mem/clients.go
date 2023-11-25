package inmemeventstream

import (
	"context"

	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/types"
)

type Clients struct {
	items map[types.UserID][]*Client
}

type Client struct {
	ctx context.Context
	ch  chan eventstream.Event

	userID  types.UserID
	eventID types.EventID
}

func (c *Clients) Add(ctx context.Context, userID types.UserID) *Client {
	client := &Client{
		ctx:     ctx,
		ch:      make(chan eventstream.Event, 16),
		eventID: types.NewEventID(),
		userID:  userID,
	}

	c.items[userID] = append(c.items[userID], client)

	return client
}

func (c *Clients) Get(userID types.UserID) []*Client {
	return c.items[userID]
}

func (c *Clients) Remove(client *Client) {
	clients, exist := c.items[client.userID]
	if !exist {
		return
	}

	for i := 0; i < len(clients); i++ {
		if clients[i].eventID == client.eventID {
			c.items[client.userID] = remove(clients, i)
			close(client.ch)
		}
	}

	if len(c.items[client.userID]) == 0 {
		delete(c.items, client.userID)
	}
}

func remove(s []*Client, i int) []*Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func NewClients() *Clients {
	return &Clients{
		items: make(map[types.UserID][]*Client),
	}
}
