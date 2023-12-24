package sendmanagermessagejob

import (
	"fmt"

	"github.com/lapitskyss/chat-service/internal/types"
)

func UnmarshalPayload(payload string) (types.MessageID, error) {
	var messageID types.MessageID
	err := messageID.UnmarshalText([]byte(payload))
	if err != nil {
		return types.MessageID{}, fmt.Errorf("unmarshal messageID: %v", err)
	}
	return messageID, nil
}

func MarshalPayload(messageID types.MessageID) (string, error) {
	if err := messageID.Validate(); err != nil {
		return "", fmt.Errorf("validate messageID: %v", err)
	}
	payload, err := messageID.MarshalText()
	if err != nil {
		return "", fmt.Errorf("marshal messageID: %v", err)
	}
	return string(payload), nil
}
