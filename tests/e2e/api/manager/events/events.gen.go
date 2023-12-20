// Package apimanagerevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package apimanagerevents

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/oapi-codegen/runtime"
)

// ChatClosedEvent defines model for ChatClosedEvent.
type ChatClosedEvent struct {
	CanTakeMoreProblems bool            `json:"canTakeMoreProblems"`
	ChatId              types.ChatID    `json:"chatId"`
	EventId             types.EventID   `json:"eventId"`
	EventType           string          `json:"eventType"`
	RequestId           types.RequestID `json:"requestId"`
}

// Event defines model for Event.
type Event struct {
	union json.RawMessage
}

// NewChatEvent defines model for NewChatEvent.
type NewChatEvent struct {
	CanTakeMoreProblems bool            `json:"canTakeMoreProblems"`
	ChatId              types.ChatID    `json:"chatId"`
	ClientId            types.UserID    `json:"clientId"`
	EventId             types.EventID   `json:"eventId"`
	EventType           string          `json:"eventType"`
	RequestId           types.RequestID `json:"requestId"`
}

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	AuthorId  types.UserID    `json:"authorId"`
	Body      string          `json:"body"`
	ChatId    types.ChatID    `json:"chatId"`
	CreatedAt time.Time       `json:"createdAt"`
	EventId   types.EventID   `json:"eventId"`
	EventType string          `json:"eventType"`
	MessageId types.MessageID `json:"messageId"`
	RequestId types.RequestID `json:"requestId"`
}

// AsNewChatEvent returns the union data inside the Event as a NewChatEvent
func (t Event) AsNewChatEvent() (NewChatEvent, error) {
	var body NewChatEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromNewChatEvent overwrites any union data inside the Event as the provided NewChatEvent
func (t *Event) FromNewChatEvent(v NewChatEvent) error {
	v.EventType = "NewChatEvent"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeNewChatEvent performs a merge with any union data inside the Event, using the provided NewChatEvent
func (t *Event) MergeNewChatEvent(v NewChatEvent) error {
	v.EventType = "NewChatEvent"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsNewMessageEvent returns the union data inside the Event as a NewMessageEvent
func (t Event) AsNewMessageEvent() (NewMessageEvent, error) {
	var body NewMessageEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromNewMessageEvent overwrites any union data inside the Event as the provided NewMessageEvent
func (t *Event) FromNewMessageEvent(v NewMessageEvent) error {
	v.EventType = "NewMessageEvent"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeNewMessageEvent performs a merge with any union data inside the Event, using the provided NewMessageEvent
func (t *Event) MergeNewMessageEvent(v NewMessageEvent) error {
	v.EventType = "NewMessageEvent"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsChatClosedEvent returns the union data inside the Event as a ChatClosedEvent
func (t Event) AsChatClosedEvent() (ChatClosedEvent, error) {
	var body ChatClosedEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromChatClosedEvent overwrites any union data inside the Event as the provided ChatClosedEvent
func (t *Event) FromChatClosedEvent(v ChatClosedEvent) error {
	v.EventType = "ChatClosedEvent"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeChatClosedEvent performs a merge with any union data inside the Event, using the provided ChatClosedEvent
func (t *Event) MergeChatClosedEvent(v ChatClosedEvent) error {
	v.EventType = "ChatClosedEvent"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t Event) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"eventType"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t Event) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "ChatClosedEvent":
		return t.AsChatClosedEvent()
	case "NewChatEvent":
		return t.AsNewChatEvent()
	case "NewMessageEvent":
		return t.AsNewMessageEvent()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t Event) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Event) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}