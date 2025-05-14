package biz

import (
	"context"
	"github.com/treewu/phicomm_dc1/app/common/data"
	data2 "github.com/treewu/phicomm_dc1/app/gateway/internal/data"
	"github.com/treewu/phicomm_dc1/pkg/xjwt"
)

type CommonBiz struct {
	wechatRepo *data2.WechatUserRepo
}

func NewCommonBiz(wechatdao *data2.WechatUserRepo) *CommonBiz {
	return &CommonBiz{
		wechatRepo: wechatdao,
	}
}

func (c *CommonBiz) GetUser(ctx context.Context) (*data2.WechatUser, error) {
	claims, err := xjwt.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	find, err := c.wechatRepo.Find(ctx, data2.WechatUser{
		UserCode: claims.UserCode,
	})
	if err != nil {
		return nil, data.UserNotLoginError("未登录")
	}
	return find, nil
}
