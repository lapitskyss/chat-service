package managerworkspace

import (
	"container/list"
	"sync"

	"github.com/lapitskyss/chat-service/internal/types"
)

type Chat struct {
	ID       types.ChatID
	ClientID types.UserID

	listItemRef *list.Element

	cursor       string
	msgMu        *sync.RWMutex
	messagesByID map[types.MessageID]*Message
	messages     []*Message
}

func NewChat(id types.ChatID, clientID types.UserID) *Chat {
	return &Chat{
		ID:           id,
		ClientID:     clientID,
		msgMu:        new(sync.RWMutex),
		messagesByID: make(map[types.MessageID]*Message),
	}
}

func (c *Chat) LastMessage() (Message, bool) {
	c.msgMu.RLock()
	defer c.msgMu.RUnlock()

	if len(c.messages) == 0 {
		return Message{}, false
	}
	return *c.messages[len(c.messages)-1], true
}

func (c *Chat) Messages() []Message {
	c.msgMu.RLock()
	defer c.msgMu.RUnlock()

	result := make([]Message, 0, len(c.messages))
	for _, m := range c.messages {
		result = append(result, *m)
	}
	return result
}

func (c *Chat) MessagesCount() int {
	c.msgMu.RLock()
	defer c.msgMu.RUnlock()

	return len(c.messages)
}

func (c *Chat) pushToFront(msg *Message) {
	c.msgMu.Lock()
	defer c.msgMu.Unlock()

	if _, ok := c.messagesByID[msg.ID]; !ok {
		c.messages = append([]*Message{msg}, c.messages...)
		c.messagesByID[msg.ID] = msg
	}
}

func (c *Chat) pushToBack(msg *Message) {
	c.msgMu.Lock()
	defer c.msgMu.Unlock()

	if _, ok := c.messagesByID[msg.ID]; !ok {
		c.messages = append(c.messages, msg)
		c.messagesByID[msg.ID] = msg
	}
}
