package biz

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	. "github.com/treewu/phicomm_dc1/app/common/data"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/data"
	"gorm.io/gorm"
	"time"
)

var planMaxCount int64 = 5

type PlanBiz struct {
	repo        *data.PlanRepo
	taskService *TaskService
	log         *log.Helper
	*CommonBiz
	parser cron.Parser
}

func NewPlanBiz(repo *data.PlanRepo, logger log.Logger, taskService *TaskService, commonBiz *CommonBiz) *PlanBiz {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	p := &PlanBiz{
		repo:        repo,
		taskService: taskService,
		log:         log.NewHelper(log.With(logger, "module", "biz/plan")),
		CommonBiz:   commonBiz,
		parser:      parser,
	}
	return p
}

func (p *PlanBiz) CreatePlan(ctx context.Context, req *v1.CreatePlanRequest) (*v1.CreatePlanReply, error) {
	user, err := p.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	if req.Plan.PlanType == v1.PlanType_PLAN_TYPE_AUTO {
		if _, err := p.parser.Parse(req.Plan.Cron); err != nil {
			return nil, ParamError("cron 解析错误")
		}
	}

	countPlan, err := p.repo.CountPlan(ctx, user.ID)
	if err != nil {
		return nil, SqlError("")
	}
	if countPlan >= planMaxCount {
		return nil, PlanExceeds(fmt.Sprintf("最多创建%d个计划", planMaxCount))
	}

	model := p.PBToModel(req.Plan)
	model.UserId = user.ID
	model.ID = 0
	plan, err := p.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	p.taskService.Add(plan)

	return &v1.CreatePlanReply{
		Id: int64(plan.ID),
	}, nil
}

func (p *PlanBiz) ListPlan(ctx context.Context, req *v1.ListPlanRequest) (*v1.ListPlanReply, error) {
	user, err := p.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	plans, err := p.repo.List(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &v1.ListPlanReply{
		Plans: slice.Map(plans, func(index int, item data.PlanModel) *v1.PlanModel {
			return p.ModelToPb(&item)
		}),
	}, nil
}

func (p *PlanBiz) SwitchPlan(ctx context.Context, req *v1.SwitchPlanRequest) (*v1.SwitchPlanReply, error) {
	plan, err := p.CheckPlanPermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	err = p.repo.SwitchPlan(ctx, req.Id, req.Enable)
	if err != nil {
		return nil, SqlError(err.Error())
	}
	plan.Enabled = req.Enable
	p.taskService.Add(plan)

	return &v1.SwitchPlanReply{}, nil
}

func (p *PlanBiz) UpdatePlan(ctx context.Context, req *v1.UpdatePlanRequest) (*v1.UpdatePlanReply, error) {
	if _, err := p.CheckPlanPermission(ctx, req.Plan.Id); err != nil {
		return nil, err
	}
	if req.Plan.PlanType == v1.PlanType_PLAN_TYPE_AUTO {
		if _, err := p.parser.Parse(req.Plan.Cron); err != nil {
			return nil, ParamError("cron 解析错误")
		}
	}
	model := p.PBToModel(req.Plan)
	err := p.repo.Update(ctx, model)
	if err != nil {
		return nil, SqlError(err.Error())
	}
	p.taskService.Add(model)
	return &v1.UpdatePlanReply{
		Id: req.Plan.Id,
	}, nil
}

func (p *PlanBiz) DeletePlan(ctx context.Context, req *v1.DeletePlanRequest) (*v1.DeletePlanReply, error) {
	plan, err := p.CheckPlanPermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	err = p.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, SqlError(err.Error())
	}
	p.taskService.Remove(plan)
	return &v1.DeletePlanReply{}, nil
}

func (p *PlanBiz) GetPlan(ctx context.Context, req *v1.GetPlanRequest) (*v1.GetPlanReply, error) {
	plan, err := p.CheckPlanPermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetPlanReply{
		Plan: p.ModelToPb(plan),
	}, nil
}

// CheckPlanPermission 检查计划权限
func (p *PlanBiz) CheckPlanPermission(ctx context.Context, id int64) (*data.PlanModel, error) {
	user, err := p.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, ParamError("id")
	}
	plan, err := p.repo.Get(ctx, id)
	if err != nil {
		return nil, RecordNotFoundError("计划不存在")
	}
	if user.ID != plan.UserId {
		return nil, ResourceWithoutPermissionError("")
	}
	return plan, nil
}

func (p *PlanBiz) ExecPlan(ctx context.Context, request *v1.ExecPlanRequest) (*v1.ExecPlanReply, error) {
	plan, err := p.CheckPlanPermission(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	p.taskService.PlanRun(plan)()
	return &v1.ExecPlanReply{
		Message: "执行成功",
	}, nil
}

func (p *PlanBiz) PBToModel(plan *v1.PlanModel) *data.PlanModel {
	d := &data.PlanModel{
		Model:    gorm.Model{ID: uint(plan.Id)},
		Name:     plan.Name,
		PlanType: int32(plan.PlanType),
		Cron:     plan.Cron,
		Enabled:  plan.Enabled,
		Devices:  nil,
	}
	// 同一个设备可能被编辑多次，这里合并一下
	with := slice.GroupWith(plan.Devices, func(device *v1.PlanDevice) string {
		return device.DeviceId
	})
	for s := range with {
		var planModelDevice = data.PlanModelDevice{}
		planModelDevice.DeviceId = s
		for _, device := range with[s] {
			planModelDevice.SwitchMain = condition.Ternary(device.SwitchMain == nil, planModelDevice.SwitchMain, device.SwitchMain)
			planModelDevice.Switch1 = condition.Ternary(device.Switch1 == nil, planModelDevice.Switch1, device.Switch1)
			planModelDevice.Switch2 = condition.Ternary(device.Switch2 == nil, planModelDevice.Switch2, device.Switch2)
			planModelDevice.Switch3 = condition.Ternary(device.Switch3 == nil, planModelDevice.Switch3, device.Switch3)
		}
		d.Devices = append(d.Devices, planModelDevice)
	}
	return d
}

func (p *PlanBiz) ModelToPb(plan *data.PlanModel) *v1.PlanModel {
	parse, _ := p.parser.Parse(plan.Cron)
	var histories []*v1.CommandHistory
	for _, commandLog := range plan.History {
		histories = append(histories, &v1.CommandHistory{
			CommandId:  commandLog.CommandId,
			ExecTime:   commandLog.CreatedAt.Format("01-02 15:04"),
			ExecResult: condition.Ternary(commandLog.Success, "成功", commandLog.Reason),
		})
	}
	return &v1.PlanModel{
		Id:       int64(plan.ID),
		Name:     plan.Name,
		PlanType: v1.PlanType(plan.PlanType),
		Cron:     plan.Cron,
		Enabled:  plan.Enabled,
		Devices: slice.Map(plan.Devices, func(index int, item data.PlanModelDevice) *v1.PlanDevice {
			return &v1.PlanDevice{
				DeviceId:   item.DeviceId,
				Switch1:    item.Switch1,
				Switch2:    item.Switch2,
				Switch3:    item.Switch3,
				SwitchMain: item.SwitchMain,
			}
		}),
		NextExecTime: parse.Next(time.Now()).Format("01-02 15:04"),
		History:      histories,
	}
}
