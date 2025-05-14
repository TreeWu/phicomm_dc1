package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/conf"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/data"
	"github.com/treewu/phicomm_dc1/pkg/middlewares"
	"github.com/treewu/phicomm_dc1/pkg/snowflake"
	"github.com/treewu/phicomm_dc1/pkg/xjwt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type WechatSession struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

type WechatBiz struct {
	cf       *conf.Data_Wechat
	repo     *data.WechatUserRepo
	log      *log.Helper
	snowNode *snowflake.Node
}

func NewWechatBiz(c *conf.Data, logger log.Logger, repo *data.WechatUserRepo) *WechatBiz {
	node, _ := snowflake.NewNode(0)

	return &WechatBiz{
		snowNode: node,
		cf:       c.Wechat,
		repo:     repo,
		log:      log.NewHelper(log.With(logger, "module", "biz/wechat")),
	}
}

// JsCode2Session https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func (b *WechatBiz) JsCode2Session(ctx context.Context, req *v1.JsCode2SessionReq) (*v1.JsCode2SessionReply, error) {
	appid := middlewares.GetMiniProgramAppid(ctx)
	client := http.Client{}
	resp, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appid, b.cf.Miniapps[appid].AppSecret, req.Code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	b.log.Infof("jscode2session resp: %s", string(all))
	var session WechatSession
	err = json.Unmarshal(all, &session)
	if err != nil {
		return nil, err
	}
	find, err := b.repo.Find(ctx, data.WechatUser{
		AppId:  appid,
		Openid: session.Openid,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			find = &data.WechatUser{
				AppId:      appid,
				Openid:     session.Openid,
				SessionKey: session.SessionKey,
				Unionid:    session.Unionid,
				CreateTime: time.Now().Unix(),
				UserCode:   "user_" + b.snowNode.Generate().Base36(),
			}
			err = b.repo.Create(ctx, *find)
			if err != nil {
				return nil, err
			}
		}
	}
	v := &v1.JsCode2SessionReply{
		Avatar:   find.Avatar,
		Nickname: find.Nickname,
		UserCode: find.UserCode,
	}

	claims := xjwt.WechatClaims{
		UserCode: find.UserCode,
		Avatar:   find.Avatar,
		Nickname: find.Nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(b.cf.TokenTtl.AsDuration())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	withClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := withClaims.SignedString([]byte(b.cf.SecretKey))
	if err != nil {
		return nil, err
	}
	v.Token = &v1.Token{AccessToken: signedString}
	return v, nil
}

func (b *WechatBiz) UpdateUser(ctx context.Context, req *v1.UpdateUserReq) (*v1.UpdateUserReply, error) {
	claims, err := xjwt.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = b.repo.Update(ctx, data.WechatUser{AppId: middlewares.GetMiniProgramAppid(ctx), UserCode: claims.UserCode, Avatar: req.Avatar, Nickname: req.Nickname})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserReply{
		Nickname: req.Avatar,
		Avatar:   req.Nickname,
	}, nil
}

func (b *WechatBiz) UserInfo(ctx context.Context, req *v1.UserInfoReq) (*v1.UserInfoReply, error) {
	claims, err := xjwt.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := b.repo.Find(ctx, data.WechatUser{AppId: middlewares.GetMiniProgramAppid(ctx), UserCode: claims.UserCode})
	if err != nil {
		return nil, err
	}
	return &v1.UserInfoReply{
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		UserCode: user.UserCode,
	}, nil
}

func (b *WechatBiz) SystemInfo(ctx context.Context, req *v1.SystemInfoReq) (*v1.SystemInfoResp, error) {
	return &v1.SystemInfoResp{
		FlushInterval: int32(b.cf.DeviceFlushInterval.AsDuration().Milliseconds()),
		Host:          b.cf.Host,
	}, nil
}

func (b *WechatBiz) CheckHost(ctx context.Context, req *v1.CheckHostReq) (*v1.CheckHostReq, error) {
	return &v1.CheckHostReq{
		Host: b.cf.Host,
	}, nil
}
