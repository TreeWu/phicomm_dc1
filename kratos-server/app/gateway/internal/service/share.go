package service

import (
	"context"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/biz"

	pb "github.com/treewu/phicomm_dc1/api/gateway/v1"
)

var _ pb.ShareHTTPServer = (*ShareService)(nil)

type ShareService struct {
	biz *biz.DeviceBiz
}

func NewShareService(biz *biz.DeviceBiz) *ShareService {
	return &ShareService{
		biz: biz,
	}
}

func (s *ShareService) ShareInvite(ctx context.Context, req *pb.ShareInviteReq) (*pb.ShareInviteReply, error) {
	return s.biz.ShareInvite(ctx, req)
}
func (s *ShareService) ShareConfirm(ctx context.Context, req *pb.ShareConfirmReq) (*pb.ShareConfirmReply, error) {
	return s.biz.ShareConfirm(ctx, req)
}

func (s *ShareService) ShareExit(ctx context.Context, req *pb.ShareExitReq) (*pb.ShareExitReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShareService) ShareRevoke(ctx context.Context, req *pb.ShareRevokeReq) (*pb.ShareRevokeReply, error) {
	return s.biz.ShareRevoke(ctx, req)

}

func (s *ShareService) GetShareList(ctx context.Context, req *pb.GetShareListReq) (*pb.GetShareListReply, error) {
	return s.biz.GetShareList(ctx, req)

}
