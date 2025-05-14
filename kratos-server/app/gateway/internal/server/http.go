package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/conf"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/service"
	"github.com/treewu/phicomm_dc1/pkg/middlewares"
)

// NewWhiteListMatcher 处于这些路由不需要验证token
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.gateway.v1.Wechat/JsCode2Session"] = struct{}{}
	whiteList["/api.gateway.v1.Wechat/SystemInfo"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	cf *conf.Data,
	commandServiceService *service.CommandServiceService,
	deviceServer *service.DeviceService,
	wechatServer *service.WechatService,
	planServer *service.PlanService,
	shareService *service.ShareService,
	logger log.Logger,
) *http.Server {
	appids := make(map[string]struct{})
	for _, app := range cf.Wechat.Miniapps {
		appids[app.AppId] = struct{}{}
	}

	var opts = []http.ServerOption{
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
		http.Middleware(
			logging.Server(logger),
			recovery.Recovery(),
			validate.Validator(),
			selector.Server(middlewares.MiniProgramAppid(appids)).Prefix("/api.gateway.v1.Wechat").Build(),
			selector.Server(
				jwt.Server(func(token *jwt5.Token) (any, error) {
					return []byte(cf.Wechat.SecretKey), nil
				}, jwt.WithSigningMethod(jwt5.SigningMethodHS256)),
			).Match(NewWhiteListMatcher()).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterCommandServiceHTTPServer(srv, commandServiceService)
	v1.RegisterDeviceHTTPServer(srv, deviceServer)
	v1.RegisterWechatHTTPServer(srv, wechatServer)
	v1.RegisterPlanHTTPServer(srv, planServer)
	v1.RegisterShareHTTPServer(srv, shareService)
	return srv
}
