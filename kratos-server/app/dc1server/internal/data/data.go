package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/treewu/phicomm_dc1/app/common/data"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, data.NewDeviceDao, NewDatabase, NewRedis, NewConsumers, data.NewUserDeviceDao)

// Data .
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewData .
func NewData(db *gorm.DB, redis *redis.Client) *Data {
	return &Data{db, redis}
}

func NewDatabase(c *conf.Data, logger log.Logger) (*gorm.DB, func(), error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})

	if err != nil {
		return nil, nil, err
	}

	_ = db.AutoMigrate(&data.Device{}, &data.DetalKWhHistory{})

	cleanup := func() {
		s, err := db.DB()
		if err != nil {
			return
		}
		s.Close()
	}
	return db, cleanup, nil
}

func NewRedis(c *conf.Data, logger log.Logger) (*redis.Client, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data/redis"))
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		helper.Errorf("redis ping error : %v", err)
	}
	helper.Infof("redis ping result : %v", pong)
	cleanup := func() {
		client.Close()
	}
	return client, cleanup, nil
}
