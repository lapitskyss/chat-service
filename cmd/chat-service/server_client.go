package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	serverclient "github.com/lapitskyss/chat-service/internal/server-client"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
)

const nameServerClient = "server-client"

func initServerClient(
	production bool,
	addr string,
	allowOrigins []string,
	kcBasePath string,
	kcRealm string,
	kcClientID string,
	kcClientSecret string,
	kcDebugMode bool,
	authResource string,
	authRole string,
	v1Swagger *openapi3.T,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	kcClient, err := keycloakclient.New(keycloakclient.NewOptions(
		kcBasePath,
		kcRealm,
		kcClientID,
		kcClientSecret,
		keycloakclient.WithDebugMode(kcDebugMode),
	))
	if err != nil {
		return nil, fmt.Errorf("create keycloack client: %v", err)
	}

	if production && kcDebugMode {
		lg.Warn("Using keycloak client debug mode in production")
	}

	srv, err := serverclient.New(serverclient.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		v1Handlers,
		kcClient,
		authResource,
		authRole,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
