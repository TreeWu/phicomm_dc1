package middlewares

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"strings"
)

type miniProgramAppidKey struct{}

type miniProgramVersionKey struct{}

func MiniProgramAppid(appids map[string]struct{}) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				referer := tr.RequestHeader().Get("referer")
				split := strings.Split(referer, "/")
				if len(split) < 5 {
					return nil, errors.New("invalid referer,appid and version not found")
				}
				if _, ok := appids[split[3]]; !ok {
					return nil, errors.New("invalid appid")
				}
				ctx = context.WithValue(ctx, miniProgramAppidKey{}, split[3])
				ctx = context.WithValue(ctx, miniProgramVersionKey{}, split[4])
			}
			return handler(ctx, req)
		}
	}
}

func GetMiniProgramAppid(ctx context.Context) string {
	return ctx.Value(miniProgramAppidKey{}).(string)
}
func GetMiniProgramVersion(ctx context.Context) string {
	return ctx.Value(miniProgramVersionKey{}).(string)
}
