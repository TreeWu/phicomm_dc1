package data

import (
	"context"
	"encoding"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"gorm.io/gorm"
	"time"
)

var _ encoding.BinaryMarshaler = (*Device)(nil)
var _ encoding.BinaryUnmarshaler = (*Device)(nil)

type Device struct {
	ID              int32  `gorm:"primary_key"`
	Mac             string `gorm:"column:mac;type:varchar(64);uniqueIndex:uniq_mac" sql:"index:device_mac" json:"mac"`
	Recover         bool   `gorm:"column:recover" json:"recover"`
	DeviceType      string `gorm:"column:device_type;type:varchar(64)" json:"device_type"`
	DeviceId        string `gorm:"column:device_id;type:varchar(64);uniqueIndex:uniq_device"  json:"device_id"`
	Key             string `json:"key" gorm:"column:key;type:varchar(64)"`
	Name            string `json:"name" gorm:"column:name;type:varchar(64)"`
	LastActivatedAt int64  `json:"last_activated_at" gorm:"column:last_activated_at;comment:上次联网激活时间"`
	LastMessageAt   int64  `json:"last_message_at" gorm:"column:last_message_at;comment:上次接收消息时间"`
	DetalKWh        int32  `json:"detal_KWh" gorm:"column:detal_KWh;comment:累计电量保存"`
	Status          *int32 `json:"status" gorm:"column:status;comment:开关状态1111"`
	I               *int32 `json:"I" gorm:"column:I;comment:电流"`
	V               *int32 `json:"V" gorm:"column:V;comment:电压"`
	P               *int32 `json:"P" gorm:"column:P;comment:功率"`
	Switch1Name     string `gorm:"column:switch1_name;type:varchar(64)" json:"switch1_name"`
	Switch2Name     string `gorm:"column:switch2_name;type:varchar(64)" json:"switch2_name"`
	Switch3Name     string `gorm:"column:switch3_name;type:varchar(64)" json:"switch3_name"`
}

func (d *Device) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Device) MarshalBinary() (data []byte, err error) {
	return json.Marshal(d)
}

type DetalKWhHistory struct {
	ID         uint   `gorm:"primary_key"`
	DeviceId   string `gorm:"column:device_id" sql:"index:device_id" json:"device_id"`
	DetalKWh   int32  `json:"detal_KWh" gorm:"column:detal_KWh;comment:累计电量累计电量"`
	CreateTime int64  `json:"create_time" gorm:"column:create_time;comment:创建时间"`
}

func (d *DetalKWhHistory) TableName() string {
	return "detal_kwh_history"
}

func (d *Device) redisKey() string {
	return "dc1:status:" + d.DeviceId
}

// TableName sets the insert table name for this struct type
func (d *Device) TableName() string {
	return "device"
}

type DeviceDao struct {
	log   *log.Helper
	db    *gorm.DB
	redis *redis.Client
}

func NewDeviceDao(db *gorm.DB, redis *redis.Client, logger log.Logger) *DeviceDao {
	return &DeviceDao{
		log:   log.NewHelper(log.With(logger, "module", "data/device")),
		db:    db,
		redis: redis,
	}
}

func (dao *DeviceDao) Insert(d *Device) error {
	return dao.db.Save(&d).Error
}

func (dao *DeviceDao) FindByMac(mac string) (*Device, error) {
	var d Device
	err := dao.db.Model(Device{Mac: mac}).Where(Device{Mac: mac}).First(&d).Error
	return &d, err
}

func (dao *DeviceDao) FindByDeviceId(deviceId string) (*Device, error) {
	var d Device
	err := dao.db.Model(&Device{}).First(&d, Device{DeviceId: deviceId}).Error
	return &d, err
}

func (dao *DeviceDao) SetDatapoint(ctx context.Context, d *Device) error {
	return dao.db.WithContext(ctx).Model(&Device{}).Where("device_id = ?", d.DeviceId).Updates(d).Error
}

// 保存设备用电量和累加总电量
func (dao *DeviceDao) DetalKWh(deviceId string, kwh int32) error {
	dao.db.Model(&deviceId).Where("device_id = ?", deviceId).Update("detal_KWh", gorm.Expr("detal_KWh + ?", kwh))
	dao.db.Create(&DetalKWhHistory{DeviceId: deviceId, DetalKWh: kwh, CreateTime: time.Now().Unix()})
	return nil
}

func (dao *DeviceDao) UpdateDevice(ctx context.Context, req *v1.UpdateDeviceRequest) error {

	updates := make(map[string]interface{})

	shouldUpdate := false
	if req.Name != nil {
		updates["name"] = *req.Name
		shouldUpdate = true
	}
	if req.Switch1Name != nil {
		updates["switch1_name"] = *req.Switch1Name
		shouldUpdate = true
	}
	if req.Switch2Name != nil {
		updates["switch2_name"] = *req.Switch2Name
		shouldUpdate = true
	}
	if req.Switch3Name != nil {
		updates["switch3_name"] = *req.Switch3Name
		shouldUpdate = true
	}
	if req.Recover != nil {
		updates["recover"] = *req.Recover
		shouldUpdate = true
	}
	if shouldUpdate {
		return dao.db.Model(&Device{}).Where("device_id = ?", *req.DeviceId).Updates(updates).Error
	}
	return nil

}

func (dao *DeviceDao) ListDevice(ctx context.Context, deviceIds []string) ([]Device, error) {

	var devices []Device
	err := dao.db.Model(&Device{}).Where("device_id in (?)", deviceIds).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}
