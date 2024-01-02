package starter

import (
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	"github.com/lapitskyss/chat-service/internal/config"
)

//nolint:unused
var clientsSet = wire.NewSet(
	provideKeycloakClient,
)

func provideKeycloakClient(cfg config.Config) (*keycloakclient.Client, error) {
	kc, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode),
	))
	if err != nil {
		return nil, fmt.Errorf("create keycloak client: %v", err)
	}
	if cfg.Global.IsProd() && cfg.Clients.Keycloak.DebugMode {
		zap.L().Warn("keycloak client in the debug mode")
	}
	return kc, nil
}
