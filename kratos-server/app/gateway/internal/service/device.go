package service

import (
	"context"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/biz"

	pb "github.com/treewu/phicomm_dc1/api/gateway/v1"
)

var _ pb.DeviceHTTPServer = (*DeviceService)(nil)

type DeviceService struct {
	biz *biz.DeviceBiz
}

func (s *DeviceService) BindDevice(ctx context.Context, req *pb.DeviceConnectReq) (*pb.DeviceConnectReply, error) {
	return s.biz.BindDevice(ctx, req)
}

func NewDeviceService(biz *biz.DeviceBiz) *DeviceService {
	return &DeviceService{
		biz: biz,
	}
}

func (s *DeviceService) UpdateDevice(ctx context.Context, req *pb.UpdateDeviceRequest) (*pb.UpdateDeviceReply, error) {
	return s.biz.UpdateDevice(ctx, req)
}

func (s *DeviceService) ListDevice(ctx context.Context, req *pb.ListDeviceRequest) (*pb.ListDeviceReply, error) {
	return s.biz.ListDevice(ctx, req)
}
