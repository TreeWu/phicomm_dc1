package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	. "github.com/treewu/phicomm_dc1/app/common/data"
	"github.com/treewu/phicomm_dc1/app/gateway/internal/data"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
	"google.golang.org/protobuf/encoding/protojson"
	"strconv"
	"time"
)

var commandFromUser = "user"

type CommandBiz struct {
	sender *data.AsyncSender
	repo   *data.CommandLogRepo
	log    *log.Helper
}

func NewCommandBiz(sender *data.AsyncSender, repo *data.CommandLogRepo, logger log.Logger) *CommandBiz {
	c := &CommandBiz{
		sender: sender,
		repo:   repo,
		log:    log.NewHelper(log.With(logger, "module", "biz/command")),
	}
	go c.reply()
	return c
}

func (c *CommandBiz) Send(ctx context.Context, req *v1.Command) (*v1.CommandReply, error) {

	commandId := strconv.FormatInt(time.Now().UnixMilli(), 10)
	req.CommandId = commandId

	protojson.Format(req)
	if err := c.sender.Send(ctx, dc1server.Dc1CommandSendQueue, dc1server.Command{
		DeviceId:    req.DeviceId,
		DeviceType:  "dc1",
		CommandId:   commandId,
		CommandTime: time.Now().UnixMilli(),
		Dc1: dc1server.CommandDc1{
			SwitchMain: req.Dc1.SwitchMain,
			Switch1:    req.Dc1.Switch_1,
			Switch2:    req.Dc1.Switch_2,
			Switch3:    req.Dc1.Switch_3,
		},
	}); err != nil {
		c.log.Warnf("send command failed: %v", err)
		return nil, ServerError("命令发送失败")
	}

	go func(ctx context.Context, req *v1.Command) {
		if err := c.repo.Create(ctx, &data.CommandLog{
			CommandId: commandId,
			DeviceId:  req.DeviceId,
			From:      commandFromUser,
			Context:   protojson.Format(req),
		}); err != nil {
			c.log.Warnf("create command log failed: %v", err)
		}
	}(ctx, req)

	return &v1.CommandReply{CommandId: commandId}, nil
}

func (c *CommandBiz) reply() {
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					c.log.Errorf("panic: %v", err)
				}
			}()
			c.sender.Reply(context.Background(), func(reply dc1server.CommandReply) {
				err := c.repo.Reply(reply)
				if err != nil {
					c.log.Warnf("reply command log failed: %v", err)
				}
			})
		}()
	}

}
