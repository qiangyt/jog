//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/qiangyt/jog/server/biz"
	"github.com/qiangyt/jog/server/conf"
	"github.com/qiangyt/jog/server/data"
	"github.com/qiangyt/jog/server/service"
)

// initServer init Jog Server with kratos application.
func initServer(ServerId, ServerName, ServerVersion, *conf.Server, *conf.Data, log.Logger) (*ServerT, func(), error) {
	panic(wire.Build(ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, NewServer))
}
