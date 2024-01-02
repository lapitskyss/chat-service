//go:build wireinject
// +build wireinject

//go:generate wire
package starter

import (
	"github.com/google/wire"
)

func Initialize() (*Service, func(), error) {
	panic(wire.Build(
		infrastructureSet,
		clientsSet,
		storageSet,
		repositorySet,
		servicesSet,
		swaggerSet,
		serverClientSet,
		serverManagerSet,
		serverDebugSet,
	))
}
