package biz

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/condition"
	"github.com/go-kratos/kratos/v2/log"
	data2 "github.com/treewu/phicomm_dc1/app/common/data"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/conf"
	"github.com/treewu/phicomm_dc1/app/dc1server/internal/data"
	. "github.com/treewu/phicomm_dc1/pkg/server/dc1server"
	"github.com/treewu/phicomm_dc1/pkg/snowflake"
	"gorm.io/gorm"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Dc1Session struct {
	DeviceId   string
	Mac        string
	DeviceType string
	*Session
}

type Dc1Pool struct {
	pool sync.Map
}

func (d *Dc1Pool) Store(deviceId string, session *Dc1Session) {
	d.pool.Store(deviceId, session)
}

func (d *Dc1Pool) Load(deviceId string) (*Dc1Session, bool) {
	load, b := d.pool.Load(deviceId)
	if b {
		return load.(*Dc1Session), true
	}
	return nil, false
}

func (d *Dc1Pool) Range(f func(string, *Dc1Session) bool) {
	d.pool.Range(func(key, value any) bool {
		return f(key.(string), value.(*Dc1Session))
	})
}

func (d *Dc1Pool) Delete(id string) {
	d.pool.Delete(id)
}

type Dc1Biz struct {
	log         *log.Helper
	Dc1Pool     Dc1Pool
	repo        *data2.DeviceDao
	sessionPool Dc1Pool
	consumer    *data.Consumers
	snowNode    *snowflake.Node
	config      *conf.Server_Dc1
	user        *data2.UserDeviceDao
}

func NewDc1Biz(repo *data2.DeviceDao, logger log.Logger, consumer *data.Consumers, c *conf.Server, user *data2.UserDeviceDao) *Dc1Biz {
	d := &Dc1Biz{
		consumer: consumer,
		repo:     repo,
		log:      log.NewHelper(log.With(logger, "module", "biz/dc1")),
		config:   c.Dc1,
		user:     user,
	}
	node, _ := snowflake.NewNode(0)
	d.snowNode = node

	go d.Heartbeat()
	go d.commandConsumers()
	go d.planCommandConsumers()
	return d
}

func (s *Dc1Biz) MessageHandler(c *Session, readString string) error {
	// 格式化打印日志
	s.log.Infof("rev msg  from [%s] : %s", c.SessionID(), readString)
	var msg UploadMessage
	err := json.Unmarshal([]byte(readString), &msg)
	if err != nil {
		s.log.Infof("format msg error : %v", err)
		c.Close()
		return err
	}
	switch {
	case msg.Action == IDENTIFY:
		// 收到 {"action":"identify","uuid":"identify2dab","auth":"cwyebqd9","params":{"device_id":"7-7170239623695255"}}
		// 回复 {"uuid":"12321323","status":200,"msg":"device identified"}
		id := msg.Params.DeviceId
		device, err := s.repo.FindByDeviceId(id)
		if err != nil {
			c.Close()
			return errors.New("设备不存在")
		}
		dc1Session := &Dc1Session{DeviceId: id, Session: c}
		s.Dc1Pool.Store(id, dc1Session)
		s.sessionPool.Store(c.SessionID(), dc1Session)
		answer := &Answer{Uuid: msg.Uuid, Status: CODE_SUCCESS, Msg: "device identified"}
		c.SendMessage(answer.ToMsg())
		if device.Recover {
			ask := NewAskWithAction(SET_DATAPOINT)
			ask.Uuid = strconv.FormatInt(time.Now().UnixMilli(), 10)
			ask.Params = AskParams{Status: condition.Ternary(device.Status == nil, 0, *device.Status)}
			c.SendMessage(ask.ToMsg())
		}
	case msg.Action == DATETIME:
		// 收到  {"action":"datetime","uuid":"datetime8011","auth":"cwyebqd9","params":{}}
		// 回复 from phicomm:  {"uuid":"datetime8011","status":200,"result":{"datetime":"2023-12-16 00:13:16"},"msg":"get datetime success"}
		answer := &Answer{Uuid: msg.Uuid, Status: CODE_SUCCESS, Result: &AnswerResult{Datetime: time.Now().Format(time.DateTime)}, Msg: "get datetime success"}
		c.SendMessage(answer.ToMsg())
	case msg.Action == ACTIVATE:
		// 判断是否激活过，如果没有，加入表，如果有，删除旧的绑定关系
		// 收到 {"action":"activate=","uuid":"activate=e28","auth":"","params":{"device_type":"PLUG_DC1_7","mac":"A4:7B:9D:06:A0:E6"}}
		// 回复 {"uuid":"activate=4c3","status":200,"msg":"device activated","result":{"uid":"7-7170239623695255","device_type":"plug","last_activated_at":"2023-12-12 23:50:36","name":"PHICOMM_","key":"cwyebqd9"}}
		var device *data2.Device
		mac := strings.ToUpper(msg.Params.Mac)
		device, err = s.repo.FindByMac(mac)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				id := s.snowNode.Generate()
				device = &data2.Device{
					Mac:             mac,
					Recover:         false,
					DeviceType:      msg.Params.DeviceType,
					DeviceId:        id.String(),
					Key:             id.Base36(),
					LastActivatedAt: time.Now().UnixMilli(),
					Name:            "PHICOMM_" + mac[9:],
					LastMessageAt:   time.Now().UnixMilli(),
				}
				err = s.repo.Insert(device)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			device.LastMessageAt = time.Now().UnixMilli()
			device.Name = "PHICOMM_" + mac[9:]
			device.Switch1Name = "开关1"
			device.Switch2Name = "开关2"
			device.Switch3Name = "开关3"
			device.LastActivatedAt = time.Now().UnixMilli()
			if err = s.repo.SetDatapoint(context.Background(), device); err != nil {
				s.log.Warnf("设备激活失败，更新激活时间,err: %v	", err)
				c.Close()
				return err
			}
		}
		if err = s.user.UnbindByDevice(context.Background(), device.DeviceId); err != nil {
			s.log.Warnf("设备激活失败，解除绑定关系,err: %v	", err)
			c.Close()
			return err
		}
		dc1Session := &Dc1Session{DeviceId: device.DeviceId, Session: c}
		s.Dc1Pool.Store(device.DeviceId, dc1Session)
		s.sessionPool.Store(dc1Session.SessionID(), dc1Session)
		answer := &Answer{Uuid: msg.Uuid, Status: CODE_SUCCESS, Msg: "device activated", Result: &AnswerResult{
			Uid:             device.DeviceId,
			DeviceType:      "plug",
			LastActivatedAt: time.UnixMilli(device.LastActivatedAt).Format(time.DateTime),
			Name:            device.Name,
			Key:             device.Key,
		}}
		c.SendMessage(answer.ToMsg())
	case msg.Action == KWH:
		// 收到:  {"action":"kWh+","uuid":"kWh+bf0ec612","auth":"cwyebqd9","params":{"detalKWh":0}}
		// 回复: {"uuid":"","result":{},"status":200,"msg":"get datetime success"}
		answer := &Answer{Uuid: msg.Uuid, Status: CODE_SUCCESS, Msg: "get datetime success", Result: &AnswerResult{}}
		c.SendMessage(answer.ToMsg())
		if value, ok := s.sessionPool.Load(c.SessionID()); ok {
			if value.DeviceId != "" && msg.Params.DetalKWh != 0 {
				err = s.repo.DetalKWh(value.DeviceId, msg.Params.DetalKWh)
				if err != nil {
					s.log.Warnf("DetalKWh error : %v", err)
				}
			}
		}
	default:
		// {"uuid":"237560166","status":200,"result":{"status":1111,"I":0,"V":235,"P":0},"msg":"get datapoint success"}
		// 应答类型消息，广播一下，让其它服务器也能收到
		if value, ok := s.sessionPool.Load(c.SessionID()); ok {
			id := value.DeviceId
			reple := CommandReply{DeviceId: id, CommandId: msg.Uuid, Success: true, ReplyTime: time.Now().UnixMilli()}
			marshal, _ := json.Marshal(reple)
			err = s.consumer.Publish(context.Background(), string(marshal))
			if err != nil {
				s.log.Warnf("publish error : %v", err)
			}
			err = s.repo.SetDatapoint(context.Background(), &data2.Device{DeviceId: id, Status: &msg.Result.Status, I: &msg.Result.I, V: &msg.Result.V, P: &msg.Result.P, LastMessageAt: time.Now().UnixMilli()})
			if err != nil {
				s.log.Warnf("SetDatapoint error : %v", err)
			}
		}

	}
	return nil
}

func (s *Dc1Biz) Heartbeat() {
	for {
		func() { // 用闭包简化 defer 作用域
			defer func() {
				if err := recover(); err != nil {
					s.log.Errorf("Heartbeat panic: %v", err)
				}
			}()
			// 核心逻辑
			msg := NewAskDataPoint().ToMsg()
			s.Dc1Pool.Range(func(key string, dc1Session *Dc1Session) bool {
				dc1Session.SendMessage(msg)
				return true
			})
		}()
		time.Sleep(s.config.HeartBeatInterval.AsDuration())
	}
}

func (s *Dc1Biz) SessionOfflineHandler(session *Session) {
	if value, ok := s.sessionPool.Load(session.SessionID()); ok {
		s.Dc1Pool.Delete(value.DeviceId)
		s.sessionPool.Delete(session.SessionID())
	}
}

func (s *Dc1Biz) SessionOnlineHandler(session *Session) {
	s.sessionPool.Store(session.SessionID(), &Dc1Session{Session: session})
}

// commandConsumers 控制命令下发，这里的开关控制命令是完整的 1111 因此可以直接下发
func (s *Dc1Biz) commandConsumers() {
	defer func() {
		if err := recover(); err != nil {
			s.log.Errorf("commandConsumers panic : %v", err)
			debug.PrintStack()
		}
		go s.commandConsumers()
	}()

	s.consumer.Consumers(context.Background(), Dc1CommandSendQueue, func(msg string) {
		var command Command
		err := json.Unmarshal([]byte(msg), &command)
		if err != nil {
			s.log.Warnf("commandConsumers error : %v", err)
		} else {
			if time.UnixMilli(command.CommandTime).Before(time.Now().Add(-s.config.CommandTimeout.AsDuration())) {
				reple := CommandReply{DeviceId: command.DeviceId, CommandId: command.CommandId, Success: false, Reason: "执行超时", ReplyTime: time.Now().UnixMilli()}
				marshal, _ := json.Marshal(reple)
				err = s.consumer.Publish(context.Background(), string(marshal))
				return
			}
			if load, ok := s.Dc1Pool.Load(command.DeviceId); ok {
				//{"action":"datapoint=","params":{"status":1111},"uuid":"864238055","auth":""}
				ask := NewAskWithAction(SET_DATAPOINT)
				ask.Uuid = command.CommandId
				ask.Params = AskParams{Status: command.Dc1.ToStatus()}
				load.SendMessage(ask.ToMsg())
			} else {
				reple := CommandReply{DeviceId: command.DeviceId, CommandId: command.CommandId, Success: false, Reason: "设备离线", ReplyTime: time.Now().UnixMilli()}
				marshal, _ := json.Marshal(reple)
				err = s.consumer.Publish(context.Background(), string(marshal))
			}
		}
	})
}

// planCommandConsumers 任务计划的命令下发
// 任务计划可能只对指定的开发进行配置，如只对 2 号开关进行配置，但是完整的协议是 1111 四位数
// 所以需要获取当前开关状态，然后替换指定开关位在进行下发
func (s *Dc1Biz) planCommandConsumers() {
	defer func() {
		if err := recover(); err != nil {
			s.log.Errorf("commandConsumers panic : %v", err)
			debug.PrintStack()
		}
		go s.planCommandConsumers()
	}()

	s.consumer.Consumers(context.Background(), Dc1CommandPlanQueue, func(msg string) {
		var command Command
		err := json.Unmarshal([]byte(msg), &command)
		if err != nil {
			s.log.Warnf("commandConsumers error : %v", err)
		} else {
			if time.UnixMilli(command.CommandTime).Before(time.Now().Add(-s.config.CommandTimeout.AsDuration())) {
				reple := CommandReply{DeviceId: command.DeviceId, CommandId: command.CommandId, Success: false, Reason: "执行超时", ReplyTime: time.Now().UnixMilli()}
				marshal, _ := json.Marshal(reple)
				err = s.consumer.Publish(context.Background(), string(marshal))
				return
			}
			if load, ok := s.Dc1Pool.Load(command.DeviceId); ok {
				device, err := s.repo.FindByDeviceId(command.DeviceId)
				if err != nil {
					s.log.Errorw("msg", "计划任务控制下发失败，设备不存在", "device", command.DeviceId)
					return
				}
				if device.Status == nil {
					device.Status = Int32Pointer(0)
				}
				dc1 := StatusToCommandDc1(*device.Status)
				dc1.SwitchMain = condition.Ternary(command.Dc1.SwitchMain != nil, command.Dc1.SwitchMain, dc1.SwitchMain)
				dc1.Switch1 = condition.Ternary(command.Dc1.Switch1 != nil, command.Dc1.Switch1, dc1.Switch1)
				dc1.Switch2 = condition.Ternary(command.Dc1.Switch2 != nil, command.Dc1.Switch2, dc1.Switch2)
				dc1.Switch3 = condition.Ternary(command.Dc1.Switch3 != nil, command.Dc1.Switch3, dc1.Switch3)
				ask := NewAskWithAction(SET_DATAPOINT)
				ask.Uuid = command.CommandId
				ask.Params = AskParams{Status: command.Dc1.ToStatus()}
				load.SendMessage(ask.ToMsg())
			} else {
				reple := CommandReply{DeviceId: command.DeviceId, CommandId: command.CommandId, Success: false, Reason: "设备离线", ReplyTime: time.Now().UnixMilli()}
				marshal, _ := json.Marshal(reple)
				err = s.consumer.Publish(context.Background(), string(marshal))
			}
		}
	})
}
