package main

import (
	"flag"
	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	"github.com/treewu/phicomm_dc1/app/dc1server/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "app/dc1server/configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, dc1 *dc1server.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			dc1,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	econf := zap.NewProductionEncoderConfig()
	econf.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(econf)
	core := zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel)
	z := zap.New(core)

	logger := log.With(kratoszap.NewLogger(z),
		"caller", log.DefaultCaller,
		//"service.id", id,
		//"service.name", Name,
		//"service.version", Version,
		//"trace.id", tracing.TraceID(),
		//"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
