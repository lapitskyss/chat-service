package inmemeventstream

import (
	"context"

	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/types"
)

type clients struct {
	items map[types.UserID][]*client
}

type client struct {
	ctx context.Context
	ch  chan eventstream.Event

	id     types.EventClientID
	userID types.UserID
}

func (c *clients) Add(ctx context.Context, userID types.UserID) *client {
	cl := &client{
		ctx: ctx,
		ch:  make(chan eventstream.Event, 1024),

		id:     types.NewEventClientID(),
		userID: userID,
	}

	c.items[userID] = append(c.items[userID], cl)

	return cl
}

func (c *clients) Get(userID types.UserID) []*client {
	return c.items[userID]
}

func (c *clients) Remove(client *client) {
	cls, exist := c.items[client.userID]
	if !exist {
		return
	}

	for i := 0; i < len(cls); i++ {
		if cls[i].id == client.id {
			c.items[client.userID] = remove(cls, i)
			close(client.ch)
		}
	}

	if len(c.items[client.userID]) == 0 {
		delete(c.items, client.userID)
	}
}

func remove(s []*client, i int) []*client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func newClients() *clients {
	return &clients{
		items: make(map[types.UserID][]*client),
	}
}
