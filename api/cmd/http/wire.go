//go:build wireinject
// +build wireinject

package http

import (
	"github.com/google/wire"
	"github.com/qxsugar/bill/api/internal"
)

func InitializeApplication() (*Application, func(), error) {
	panic(wire.Build(
		NewApplication,
		internal.MiscProviderSet,
		internal.DaoProviderSet,
		internal.ServiceProviderSet,
		internal.RouterProviderSet,
	))
}
