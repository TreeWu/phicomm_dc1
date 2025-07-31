package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"github.com/gorilla/handlers"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer() *http.Server {

	var opts = []http.ServerOption{
		http.Address(":8081"),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
		http.Middleware(
			recovery.Recovery(),
		),
	}
	srv := http.NewServer(opts...)
	srv.Handle("/debug/pprof/", pprof.NewHandler())

	return srv
}
