package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/biz"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/conf"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
)

func NewDc1Server(c *conf.Server, logger log.Logger, dc1 *biz.Dc1Biz) *dc1server.Server {

	var opts = []dc1server.ServerOption{}

	if c.Dc1.Network != "" {
		opts = append(opts, dc1server.WithNetwork(c.Dc1.Network))
	}
	if c.Dc1.Addr != "" {
		opts = append(opts, dc1server.WithAddress(c.Dc1.Addr))
	}

	opts = append(opts, dc1server.WithLogger(logger),
		dc1server.WithMessageHandler(dc1.MessageHandler),
		dc1server.WithSessionOfflineHandler(dc1.SessionOfflineHandler),
		dc1server.WithSessionOnlineHandler(dc1.SessionOnlineHandler),
	)

	server := dc1server.NewServer(opts...)

	return server
}
