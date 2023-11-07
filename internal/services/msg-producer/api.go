package msgproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/lapitskyss/chat-service/internal/types"
)

type Message struct {
	ID         types.MessageID `json:"id"`
	ChatID     types.ChatID    `json:"chatId"`
	Body       string          `json:"body"`
	FromClient bool            `json:"fromClient"`
}

func (s *Service) ProduceMessage(ctx context.Context, msg Message) error {
	key, err := msg.ChatID.MarshalText()
	if err != nil {
		return fmt.Errorf("marshal chat id: %v", err)
	}

	value, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %v", err)
	}

	if s.cipher != nil {
		nonce, err := s.nonceFactory(s.cipher.NonceSize())
		if err != nil {
			return fmt.Errorf("genegate nonce: %v", err)
		}
		value = s.cipher.Seal(nil, nonce, value, nil)
		value = append(nonce, value...)
	}

	err = s.wr.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("write messages: %v", err)
	}
	return nil
}

func (s *Service) Close() error {
	return s.wr.Close()
}
