package service

import (
	"context"
	pb "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/biz"
)

var _ pb.WechatHTTPServer = (*WechatService)(nil)

type WechatService struct {
	biz *biz.WechatBiz
}

func (s *WechatService) CheckHost(ctx context.Context, req *pb.CheckHostReq) (*pb.CheckHostReq, error) {
	return s.biz.CheckHost(ctx, req)
}

func (s *WechatService) SystemInfo(ctx context.Context, req *pb.SystemInfoReq) (*pb.SystemInfoResp, error) {
	return s.biz.SystemInfo(ctx, req)

}

func (s *WechatService) UserInfo(ctx context.Context, req *pb.UserInfoReq) (*pb.UserInfoReply, error) {
	return s.biz.UserInfo(ctx, req)
}

func NewWechatService(biz *biz.WechatBiz) *WechatService {
	return &WechatService{
		biz: biz,
	}
}

func (s *WechatService) JsCode2Session(ctx context.Context, req *pb.JsCode2SessionReq) (*pb.JsCode2SessionReply, error) {
	return s.biz.JsCode2Session(ctx, req)
}
func (s *WechatService) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserReply, error) {
	return s.biz.UpdateUser(ctx, req)
}
