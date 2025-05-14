package dc1server

import (
	"encoding/json"
	"strconv"
	"time"
)

type UploadMessage struct {
	Action string `json:"action"`
	Uuid   string `json:"uuid"`
	Auth   string `json:"auth"`
	Status int32  `json:"status"`
	Msg    string `json:"msg"`
	Result struct {
		Status int32 `json:"status"` // 1110 开关状态
		I      int32 `json:"I"`      //电流
		V      int32 `json:"V"`      // 电压
		P      int32 `json:"P"`      // 功率
	}
	Params struct {
		DeviceId   string `json:"device_id"`
		Mac        string `json:"mac"`
		DeviceType string `json:"device_type"`
		DetalKWh   int32  `json:"detalKWh"`
	} `json:"params"`
}

type Answer struct {
	Uuid   string        `json:"uuid"`
	Status int           `json:"status"`
	Msg    string        `json:"msg"`
	Result *AnswerResult `json:"result,omitempty"`
}

type AnswerResult struct {
	Uid             string `json:"uid"`
	DeviceType      string `json:"json_type"`
	Plug            string `json:"plug"`
	LastActivatedAt string `json:"last_activated_at"`
	Name            string `json:"name"`
	Key             string `json:"key"`
	Datetime        string `json:"datetime"`
}

func (a Answer) ToMsg() []byte {
	msg, _ := json.Marshal(a)
	return append(msg, '\n')
}

type Ask struct {
	Action string    `json:"action"`
	Uuid   string    `json:"uuid"`
	Auth   string    `json:"auth"`
	Params AskParams `json:"params"`
}
type AskParams struct {
	Status int32 `json:"status"` // 1110 开关状态
}

func NewAskDataPoint() Ask {
	return Ask{
		Action: DATAPOINT,
		Uuid:   strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}
func NewAskWithAction(action string) Ask {
	return Ask{
		Action: action,
		Uuid:   strconv.FormatInt(time.Now().UnixMilli(), 10),
	}
}
func (a Ask) ToMsg() []byte {
	msg, _ := json.Marshal(a)
	return append(msg, '\n')
}

type Command struct {
	DeviceId    string     `json:"device_id"`
	DeviceType  string     `json:"device_type"`
	CommandId   string     `json:"command_id"`
	CommandTime int64      `json:"command_time"`
	Dc1         CommandDc1 `json:"dc1"`
}

type CommandDc1 struct {
	SwitchMain *int32 `json:"switch_main"`
	Switch1    *int32 `json:"switch_1"`
	Switch2    *int32 `json:"switch_2"`
	Switch3    *int32 `json:"switch_3"`
}

func StatusToCommandDc1(status int32) CommandDc1 {
	dc1 := CommandDc1{
		Switch3:    Int32Pointer(status / 1000),
		Switch2:    Int32Pointer(status % 1000 / 100),
		Switch1:    Int32Pointer(status % 1000 % 100 / 10),
		SwitchMain: Int32Pointer(status % 1000 % 100 % 10),
	}
	return dc1
}

// ToStatus 转换为状态
// 可以在这里设置是否分开总开
func (c CommandDc1) ToStatus() int32 {
	var result int32
	if c.SwitchMain != nil {
		result += *c.SwitchMain
	}
	if c.Switch1 != nil {
		result += *c.Switch1 * 10
	}
	if c.Switch2 != nil {
		result += *c.Switch2 * 100
	}
	if c.Switch3 != nil {
		result += *c.Switch3 * 1000
	}
	return result
}

type CommandReply struct {
	CommandId string `json:"command_id"`
	DeviceId  string `json:"device_id"`
	Success   bool   `json:"result"`
	Reason    string `json:"reason"`
	ReplyTime int64  `json:"reply_time"`
}

func Int32Pointer(i int32) *int32 {
	return &i
}
