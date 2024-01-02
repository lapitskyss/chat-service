package starter

import (
	"fmt"

	"github.com/google/wire"

	"github.com/lapitskyss/chat-service/internal/config"
	serverdebug "github.com/lapitskyss/chat-service/internal/server-debug"
)

//nolint:unused
var serverDebugSet = wire.NewSet(
	provideServerDebug,
)

func provideServerDebug(
	cfg config.Config,
	clientV1Swagger ClientV1Swagger,
	managerV1Swagger ManagerV1Swagger,
	clientEventsSwagger ClientEventsSwagger,
	managerEventsSwagger ManagerEventsSwagger,
) (*serverdebug.Server, error) {
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(
		cfg.Servers.Debug.Addr,
		clientV1Swagger,
		managerV1Swagger,
		clientEventsSwagger,
		managerEventsSwagger,
	))
	if err != nil {
		return nil, fmt.Errorf("init debug server: %v", err)
	}
	return srvDebug, nil
}
