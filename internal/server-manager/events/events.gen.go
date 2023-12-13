// Package managerevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package managerevents

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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/oapi-codegen/runtime"
)

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
	case "NewChatEvent":
		return t.AsNewChatEvent()
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

	"H4sIAAAAAAAC/7SUMW8bPQyG/8qB3zeeLQddAo1tOmRIWrTpFGSQdbRPtSSqlO5cw/B/LyhfYrfoUBjw",
	"ZIGk+D6vSd0eLIVEEWPJoPeQbY/B1OPHEWORQ+eyZRdcNIVYAsGk5OJajo+4/dCbMtXCf+rUTk291G81",
	"hxYSU0Iuu0cTEDSgxJ92CSVHET+tQD/v4X/G1b91fDm0f3DsXzUcVivWxCezwQdi/My09BhquIiohiWR",
	"RxNF3vam3HeSWxEHI5aGwXXQvtbmwuK8hZ+zNc2moPzkuejf352nZi4k4iOPKT1oWLvSD8u5paC8Sa7k",
	"zS5nJaqzjDw6i8rFghyNV7UrHITKO4wXc33LyNfhqpO7FKuO6opcdaNOQ55IDi0w/hgwX8z9Zbp+BfIJ",
	"zjF2oJ//urVvK3q2FadBnFs/N/ryZoyW39GWo5aLK6r/kCtecu9N3DRfhyROGtnm5sFEs0Zu6qxEfETO",
	"jiJoGG/qc00YTXKg4d38Zr6AttrPoOPgfQviFDnX99yhfEVSOV6/wxE9pYCxNMcqaGFgDxq2WSvlyRrf",
	"Uy76dnG7UNssz/xXAAAA//8WwqQprQQAAA==",
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
