package service

import (
	"context"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/biz"

	pb "github.com/treewu/phicomm_dc1/api/gateway/v1"
)

type CommandServiceService struct {
	command *biz.CommandBiz
	pb.UnimplementedCommandServiceServer
}

func NewCommandServiceService(command *biz.CommandBiz) *CommandServiceService {
	return &CommandServiceService{
		command: command,
	}
}

func (s *CommandServiceService) SendCommand(ctx context.Context, req *pb.Command) (*pb.CommandReply, error) {
	return s.command.Send(ctx, req)
}
