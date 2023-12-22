package starter

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/wire"

	clientevents "github.com/lapitskyss/chat-service/internal/server-client/events"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	managerevents "github.com/lapitskyss/chat-service/internal/server-manager/events"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
)

type (
	ClientV1Swagger      *openapi3.T
	ManagerV1Swagger     *openapi3.T
	ClientEventsSwagger  *openapi3.T
	ManagerEventsSwagger *openapi3.T
)

//nolint:unused
var swaggerSet = wire.NewSet(
	provideClientV1Swagger,
	provideManagerV1Swagger,
	provideClientEventsSwagger,
	provideManagerEventsSwagger,
)

func provideClientV1Swagger() (ClientV1Swagger, error) {
	clientV1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get client v1 swagger: %v", err)
	}
	return clientV1Swagger, nil
}

func provideManagerV1Swagger() (ManagerV1Swagger, error) {
	managerV1Swagger, err := managerv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get manager v1 swagger: %v", err)
	}
	return managerV1Swagger, nil
}

func provideClientEventsSwagger() (ClientEventsSwagger, error) {
	clientEventsSwagger, err := clientevents.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get client events swagger: %v", err)
	}
	return clientEventsSwagger, nil
}

func provideManagerEventsSwagger() (ManagerEventsSwagger, error) {
	managerEventsSwagger, err := managerevents.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get manager events swagger: %v", err)
	}
	return managerEventsSwagger, nil
}
