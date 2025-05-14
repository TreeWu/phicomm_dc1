package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
	"gorm.io/gorm"
)

type CommandLog struct {
	gorm.Model
	DeviceId  string `gorm:"column:device_id;type:varchar(64);index:command_id" json:"device_id"`
	CommandId string `gorm:"column:command_id;type:varchar(64);index:command_id" json:"command_id"`
	PlanId    uint   `gorm:"column:plan_id;type:int;index:plan_id;default:0" json:"plan_id"`
	From      string `gorm:"column:from;type:varchar(32)" json:"from"`
	Context   string `gorm:"column:context" json:"context"`
	ReplyTime int64  `gorm:"column:reply_time" json:"reply_time"`
	Success   bool   `gorm:"column:success" json:"success"`
	Reason    string `gorm:"column:reason;type:varchar(64)" json:"reason"`
}

type CommandLogRepo struct {
	db  *gorm.DB
	log *log.Helper
}

func NewCommandLogRepo(data *Data, logger log.Logger) *CommandLogRepo {
	return &CommandLogRepo{
		db:  data.db,
		log: log.NewHelper(log.With(logger, "module", "data/commandLog")),
	}
}

func (c *CommandLogRepo) Create(ctx context.Context, commandLog *CommandLog) error {
	return c.db.Create(commandLog).Error
}

func (c *CommandLogRepo) Reply(reply dc1server.CommandReply) error {
	return c.db.Model(&CommandLog{}).Where("command_id = ? and device_id = ?", reply.CommandId, reply.DeviceId).Updates(CommandLog{
		ReplyTime: reply.ReplyTime,
		Success:   reply.Success,
		Reason:    reply.Reason,
	}).Error

}
