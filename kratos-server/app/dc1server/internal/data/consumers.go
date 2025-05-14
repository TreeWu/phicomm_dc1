package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/treewu/phicomm_dc1/pkg/server/dc1server"
	"time"
)

type Consumers struct {
	redis *redis.Client
	log   *log.Helper
}

func NewConsumers(redis *redis.Client, logger log.Logger) *Consumers {
	return &Consumers{
		log:   log.NewHelper(log.With(logger, "module", "data/consumers")),
		redis: redis,
	}
}

func (c *Consumers) Consumers(ctx context.Context, queue string, f func(string)) {
	for {

		result, err := c.redis.BRPop(ctx, time.Minute, queue).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			c.log.Errorf("redis brpop error: %v", err)
			return
		}
		f(result[1])
	}
}

func (c *Consumers) Publish(ctx context.Context, msg string) error {
	return c.redis.Publish(ctx, dc1server.Dc1StatusReplyQueue, msg).Err()
}
