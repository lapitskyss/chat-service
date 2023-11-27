// Package clientevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package clientevents

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/oapi-codegen/runtime"
)

// Event defines model for Event.
type Event struct {
	union json.RawMessage
}

// EventCommon defines model for EventCommon.
type EventCommon struct {
	EventId   types.EventID   `json:"eventId"`
	EventType string          `json:"eventType"`
	MessageId types.MessageID `json:"messageId"`
	RequestId types.RequestID `json:"requestId"`
}

// MessageSentEvent defines model for MessageSentEvent.
type MessageSentEvent = EventCommon

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	AuthorId  *types.UserID   `json:"authorId,omitempty"`
	Body      string          `json:"body"`
	CreatedAt time.Time       `json:"createdAt"`
	EventId   types.EventID   `json:"eventId"`
	EventType string          `json:"eventType"`
	IsService bool            `json:"isService"`
	MessageId types.MessageID `json:"messageId"`
	RequestId types.RequestID `json:"requestId"`
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

// AsMessageSentEvent returns the union data inside the Event as a MessageSentEvent
func (t Event) AsMessageSentEvent() (MessageSentEvent, error) {
	var body MessageSentEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromMessageSentEvent overwrites any union data inside the Event as the provided MessageSentEvent
func (t *Event) FromMessageSentEvent(v MessageSentEvent) error {
	v.EventType = "MessageSentEvent"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeMessageSentEvent performs a merge with any union data inside the Event, using the provided MessageSentEvent
func (t *Event) MergeMessageSentEvent(v MessageSentEvent) error {
	v.EventType = "MessageSentEvent"
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
	case "MessageSentEvent":
		return t.AsMessageSentEvent()
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

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RUT2/bPgz9KgZ/P2AXJ063S+Hb1g5DMawFlu1U9KDYjK3WFjWRThYE/u4DZa9zk6Er",
	"AvQkguKfp8cn7qGg1pNDJwz5HriosTXR/LhBJ2qUlotgW+uMUFBHa7y3rlLzCzKbCpfoZIyH/7I/JbOx",
	"XnYUl8I1bkfvs5mHYX0KPpDHILtr0yLkgOr/tvOod+TwZg357R7+D7h+edHn44/g93d9OjB0QW1LTrkY",
	"YVmM9EVUV6Waawqt0Qd2nS0hBVGsObAEZTGFn7OKZqNTD57HyleX07uZbT2FOBBvpIYcKit1t5oX1GaN",
	"8Vb4YcecFbWRGWPY2AIz6wSDM00Wy0LfpxO28v0Bkj6FdnjoqbhHnl4HecAfHfLJnH4d018B2wjOBiwh",
	"v30c/ZTsKbXTp9w9QqfVPRZR38efag+maV6g66kio0SPPtlJldJDcZtOagqnTuI7Y3gdiayo3P1V10VA",
	"I1i+lyeISyM4E9viEew+BcvLodGk4IqoQePgcOKx77TLNP14xHdRMtatKda20ujtB+MekmXnlYnkojaS",
	"XDQWnSRxFgwpbDCw1WUDm7O47Dw64y3k8G5+Nl9AGtmLE8pYupUaFQ5LHHWJexnSryTpGDlZU0gqdBiM",
	"WFclUa48T26kxrC1jImVpCRk90bmEPtpJDmdPHxCWWoTpYI9OR608Xax0KMgJ78V531ji5iY3fOwLAed",
	"qfVPFQ4f7OkDbj6rV/2qBgwcBf005hI32JBvlcIhClLoQgM5bDnPsoYK09TEkp8vzhfZlnUwvwIAAP//",
	"9Qfe+REHAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}