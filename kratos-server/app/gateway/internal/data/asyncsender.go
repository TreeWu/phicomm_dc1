package data

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
)

type AsyncSender struct {
	redisClient *redis.Client
	log         *log.Helper
}

func (a AsyncSender) Reply(ctx context.Context, f func(reply dc1server.CommandReply)) error {
	subscribe := a.redisClient.Subscribe(ctx, dc1server.Dc1StatusReplyQueue)
	// 检查订阅是否成功
	_, err := subscribe.Receive(ctx)
	if err != nil {
		a.log.Warnf("Error subscribing to channel: %v", err)
	}

	for {
		select {
		case msg := <-subscribe.Channel():
			var c dc1server.CommandReply
			err = json.Unmarshal([]byte(msg.Payload), &c)
			if err != nil {
				a.log.Warnf("Error unmarshalling message: %v", err)
			} else {
				f(c)
			}
		}
	}
}

func (a AsyncSender) Send(ctx context.Context, queue string, req dc1server.Command) error {
	msg, _ := json.Marshal(req)
	return a.redisClient.LPush(ctx, queue, string(msg)).Err()
}

func NewAsyncSender(data *Data, logger log.Logger) *AsyncSender {
	return &AsyncSender{redisClient: data.redis, log: log.NewHelper(log.With(logger, "module", "data/asyncSender"))}
}
