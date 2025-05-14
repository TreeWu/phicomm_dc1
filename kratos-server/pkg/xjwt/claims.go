package xjwt

import (
	"context"
	jwt2 "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
)

type WechatClaims struct {
	UserCode string `json:"userCode"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

func FromContext(ctx context.Context) (*WechatClaims, error) {
	token, ok := jwt2.FromContext(ctx)
	if !ok {
		return nil, jwt2.ErrMissingJwtToken
	}

	claims := token.(jwt.MapClaims)

	return &WechatClaims{
		UserCode: claims["userCode"].(string),
		Avatar:   claims["avatar"].(string),
		Nickname: claims["nickname"].(string),
	}, nil
}
