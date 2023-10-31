package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func Encode(data any) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("json marshal: %v", err)
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func Decode(in string, to any) error {
	data, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return fmt.Errorf("base64 decode string: %v", err)
	}
	err = json.Unmarshal(data, &to)
	if err != nil {
		return fmt.Errorf("json unmarhsal: %v", err)
	}
	return nil
}
