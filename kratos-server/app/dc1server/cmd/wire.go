//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/biz"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/conf"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/data"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/server"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, newApp, biz.ProviderSet, data.ProviderSet))
}
