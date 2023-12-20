package managerevents

import (
	"encoding/json"
	"fmt"
	"io"
)

func Decode(r io.Reader) (any, error) {
	var dst ReadEvent
	err := json.NewDecoder(r).Decode(&dst)
	if err != nil {
		return nil, fmt.Errorf("unmarshal read event, %v", err)
	}
	v, err := dst.ValueByDiscriminator()
	if err != nil {
		return nil, fmt.Errorf("unmarshal read event, %v", err)
	}
	return v, nil
}
