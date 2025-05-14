package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/data"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
	"strconv"
	"time"
)

type TaskService struct {
	cronService *cron.Cron
	sender      *data.AsyncSender
	log         *log.Helper
	cronMap     map[uint]cron.EntryID
	planRepo    *data.PlanRepo
	commandLog  *data.CommandLogRepo
}

func NewTaskService(sender *data.AsyncSender, repo *data.PlanRepo, commandLog *data.CommandLogRepo, logger log.Logger) (*TaskService, func()) {
	c := cron.New(cron.WithSeconds(), cron.WithParser(cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.DowOptional|cron.Descriptor)))
	t := &TaskService{
		cronService: c,
		sender:      sender,
		log:         log.NewHelper(log.With(logger, "module", "biz/plan")),
		cronMap:     make(map[uint]cron.EntryID),
		planRepo:    repo,
		commandLog:  commandLog,
	}
	t.cronService.Start()
	t.initialize()
	return t, func() {
		t.cronService.Stop()
	}
}

func (c *TaskService) initialize() {
	enables, err := c.planRepo.FindEnable(context.Background())
	if err != nil {
		c.log.Errorw("msg", "list enable plan failed", "err", err)
		return
	}
	for _, enable := range enables {
		c.Add(&enable)
	}
}
func (c *TaskService) Remove(d *data.PlanModel) {
	if _, ok := c.cronMap[d.ID]; ok {
		c.cronService.Remove(c.cronMap[d.ID])
	}
}

// Add 添加计划任务
// 先通过ID删除当前计划，再判断是否需要启动新的计划
// 因为在 update 或者 switch 的时候，计划类型可能由 自动变成手动 或者 启用变停止
func (c *TaskService) Add(p *data.PlanModel) {
	c.Remove(p)
	if p.PlanType != int32(v1.PlanType_PLAN_TYPE_AUTO) || p.Enabled == false {
		return
	}
	addFunc, err := c.cronService.AddFunc(p.Cron, c.PlanRun(p))
	if err != nil {
		c.log.Errorw("msg", "add cron failed", "err", err)
	}
	c.cronMap[p.ID] = addFunc
}

func (c *TaskService) PlanRun(p *data.PlanModel) func() {
	return func() {
		commandId := strconv.FormatInt(time.Now().UnixMilli(), 10)
		for _, device := range p.Devices {
			command := dc1server.Command{
				DeviceId:    device.DeviceId,
				CommandId:   commandId,
				CommandTime: time.Now().UnixMilli(),
				Dc1: dc1server.CommandDc1{
					SwitchMain: device.SwitchMain,
					Switch1:    device.Switch1,
					Switch2:    device.Switch2,
					Switch3:    device.Switch3,
				},
			}
			err := c.sender.Send(context.Background(), dc1server.Dc1CommandPlanQueue, command)
			if err != nil {
				c.log.Errorw("msg", "send command failed", "err", err)
			}
			if err := c.commandLog.Create(context.Background(), &data.CommandLog{
				DeviceId:  command.DeviceId,
				CommandId: command.CommandId,
				From:      "plan_auto",
				PlanId:    p.ID,
				Context:   strconv.Itoa(int(command.Dc1.ToStatus())),
			}); err != nil {
				c.log.Errorw("msg", "create command log failed", "err", err)
			}

		}
	}
}
