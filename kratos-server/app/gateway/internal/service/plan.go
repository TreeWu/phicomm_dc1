package service

import (
	"context"
	pb "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/biz"
)

var _ pb.PlanHTTPServer = (*PlanService)(nil)

type PlanService struct {
	plan *biz.PlanBiz
}

func (s *PlanService) ExecPlan(ctx context.Context, request *pb.ExecPlanRequest) (*pb.ExecPlanReply, error) {
	return s.plan.ExecPlan(ctx, request)
}

func NewPlanService(plan *biz.PlanBiz) *PlanService {
	return &PlanService{
		plan: plan,
	}
}

func (s *PlanService) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.CreatePlanReply, error) {
	return s.plan.CreatePlan(ctx, req)
}
func (s *PlanService) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.UpdatePlanReply, error) {
	return s.plan.UpdatePlan(ctx, req)
}
func (s *PlanService) DeletePlan(ctx context.Context, req *pb.DeletePlanRequest) (*pb.DeletePlanReply, error) {
	return s.plan.DeletePlan(ctx, req)
}
func (s *PlanService) GetPlan(ctx context.Context, req *pb.GetPlanRequest) (*pb.GetPlanReply, error) {
	return s.plan.GetPlan(ctx, req)
}
func (s *PlanService) ListPlan(ctx context.Context, req *pb.ListPlanRequest) (*pb.ListPlanReply, error) {
	return s.plan.ListPlan(ctx, req)
}

func (s *PlanService) SwitchPlan(ctx context.Context, request *pb.SwitchPlanRequest) (*pb.SwitchPlanReply, error) {
	return s.plan.SwitchPlan(ctx, request)
}
